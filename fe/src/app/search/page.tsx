"use client";
import React, { useState, useEffect } from "react";
import Image from "next/image";
import { source, destination, dummyData } from "./data";
import Graph from "../components/chart";
import { TabsDemo } from "../components/tabs/tabs";

export default function Search() {
  const [isRotating, setIsRotating] = useState(false);
  const [selectedAlgorithm, setSelectedAlgorithm] = useState("BFS");
  const [showResult, setShowResult] = useState(false);

  const [currentSourceIndex, setCurrentSourceIndex] = useState(0);
  const [currentDestinationIndex, setCurrentDestinationIndex] = useState(0);

  const [sourceInput, setSourceInput] = useState("");
  const [sourceResults, setSourceResults] = useState([]);
  const [sourceFocused, setSourceFocused] = useState(false);

  const [destinationInput, setDestinationInput] = useState("");
  const [destinationResults, setDestinationResults] = useState([]);
  const [destinationFocused, setDestinationFocused] = useState(false);

  const [paths, setPaths] = useState([]);
  const [degrees, setDegrees] = useState(0);
  const [executionTime, setExecutionTime] = useState(0);

  const fetchResults = async (input: string, setResult: any) => {
    if (input) {
      try {
        const response = await fetch(
          `https://en.wikipedia.org/w/api.php?action=opensearch&search=${input}&format=json&origin=*`
        );
        const data = await response.json();
        setResult(data[1]);
      } catch (error) {
        console.error("Failed to fetch data:", error);
      }
    } else {
      setResult([]);
    }
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      fetchResults(sourceInput, setSourceResults);
    }, 500);
    console.log(sourceInput);
    return () => clearTimeout(timer);
  }, [sourceInput]);

  useEffect(() => {
    const timer = setTimeout(() => {
      fetchResults(destinationInput, setDestinationResults);
    }, 500);
    console.log(destinationInput);
    return () => clearTimeout(timer);
  }, [destinationInput]);

  useEffect(() => {
    const updatePlaceholders = () => {
      setCurrentSourceIndex((prevIndex) => (prevIndex + 1) % source.length);

      setCurrentDestinationIndex(
        (prevIndex) => (prevIndex + 1) % destination.length
      );
    };

    const intervalId = setInterval(updatePlaceholders, 2000);

    return () => clearInterval(intervalId);
  }, []);

  const handleAlgorithmSelection = (algorithm: string) => {
    setSelectedAlgorithm(algorithm);
    setShowResult(false);
  };

  const handleSearch = async () => {
    if (!selectedAlgorithm || !sourceInput || !destinationInput) {
      alert(
        "Please select an algorithm and enter both a source and a destination."
      );
      return;
    }

    try {
      const response = await fetch(`http://localhost:4000/search/BFS`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ start: sourceInput, goal: destinationInput }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      console.log("Data received:", data);

      // Update state with the response data
      setShowResult(true);
      setPaths(data.paths.length);
      setDegrees(data.degrees);
      setExecutionTime(data.executionTime);
    } catch (error) {
      console.error("There was an error!", error);
    }
  };

  const swapInputs = () => {
    setIsRotating(true);
    setTimeout(() => { 
      setIsRotating(false);
    }, 500);

    setSourceInput(destinationInput);
    setDestinationInput(sourceInput);
    setCurrentSourceIndex(currentDestinationIndex);
    setCurrentDestinationIndex(currentSourceIndex);
  };

  return (
    <div className="max-w-3xl w-full h-full mt-8">
      <div className="">
        <div className="bg-white p-8 rounded-xl">
          <div className="flex flex-col items-start mb-8">
            <div className="text-3xl sm:text-4xl text-[#075A5A] font-bold">
              Search Algorithm
            </div>
            <div className="mt-4 mb-4 flex flex-col">
              <TabsDemo onAlgorithmChange={handleAlgorithmSelection} />
            </div>
            <div className="h-[1px] mt-4 w-full bg-[#A1B6BA]"></div>
          </div>
          <div className="flex flex-col bg-white ">
            <div className="bg-[#E4F8F8] relative w-full h-[120px] py-4 rounded-lg px-8">
              <div className="text-[#3F6870] text-3xl flex flex-row justify-items-start">
                From
              </div>
              <input
                className="absolute bg-[#E4F8F8] bottom-4 left-8 right-0 focus:outline-none text-[#075A5A] font-black text-xl sm:text-3xl"
                type="text"
                placeholder={source[currentSourceIndex]}
                value={sourceInput}
                onChange={(e) => setSourceInput(e.target.value)}
                onFocus={() => setSourceFocused(true)}
                onBlur={() => setTimeout(() => setSourceFocused(false), 100)}
              />
              {sourceResults.length > 0 && sourceFocused && (
                <ul className="absolute max-h-72 overflow-y-auto bg-white w-full left-0 right-0 top-32 z-10 border border-gray-300 rounded-lg px-6 py-2">
                  {sourceResults.map((result, index) => (
                    <li
                      key={index}
                      className="p-2 hover:bg-gray-100 rounded-lg cursor-pointer"
                      onMouseDown={(e) => {
                        e.preventDefault(); // Prevent the input from losing focus
                      }}
                      onClick={() => {
                        setSourceInput(result);
                        setSourceResults([]);
                        setSourceFocused(false);
                      }}
                    >
                      {result}
                    </li>
                  ))}
                </ul>
              )}
            </div>
            <button
              className={`flex justify-center mt-6 mb-6 items-center flex-cold ${
                isRotating ? "rotate-animation" : ""
              } `}
              onClick={swapInputs}
            >
              <Image src="/arrow.svg" alt="arrow" width={30} height={30} />
            </button>
            <div className="bg-[#E4F8F8] relative w-full h-[120px] rounded-lg py-4 px-8">
              <div className="text-[#3F6870] text-3xl flex flex-row justify-items-start">
                To
              </div>
              <input
                className="absolute bg-[#E4F8F8] bottom-4 left-8 right-0 focus:outline-none text-[#075A5A] font-black text-xl sm:text-3xl"
                type="text"
                placeholder={destination[currentDestinationIndex]}
                value={destinationInput}
                onChange={(e) => setDestinationInput(e.target.value)}
                onFocus={() => setDestinationFocused(true)}
                onBlur={() =>
                  setTimeout(() => setDestinationFocused(false), 100)
                }
              />
              {destinationResults.length > 0 && destinationFocused && (
                <ul className="absolute max-h-72 overflow-y-auto bg-white w-full left-0 right-0 top-32 z-10 border border-gray-300 rounded-lg px-6 py-2">
                  {destinationResults.map((result, index) => (
                    <li
                      key={index}
                      className="p-2 hover:bg-gray-100 rounded-lg cursor-pointer"
                      onMouseDown={(e) => {
                        e.preventDefault(); // Prevent the input from losing focus
                      }}
                      onClick={() => {
                        setDestinationInput(result);
                        setDestinationResults([]);
                        setDestinationFocused(false);
                      }}
                    >
                      {result}
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>
          <div className="flex flex-row justify-center items-center mt-8">
            <button
              className="bg-[#075A5A] font-bold text-white w-full h-[72px] rounded-lg border border-2 border-[#1A535C] text-2xl sm:text-3xl font-bold"
              onClick={handleSearch}
            >
              Go!
            </button>
          </div>

          {showResult && (
            <div className="mt-12 result flex flex-col items-center justify-center">
              <div className="w-full max-w-5xl mt-8">
                <div className="text-3xl text-center text-[#1A535C]">
                  Found <strong>{paths.length} path(s)</strong> with{" "}
                  <strong>{degrees} degrees</strong> of separation from{" "}
                  <strong>{sourceInput}</strong> to{" "}
                  <strong>{destinationInput}</strong> in{" "}
                  <strong>{executionTime} seconds</strong> using{" "}
                  <strong>{selectedAlgorithm} Algorithm.</strong>
                </div>
              </div>
              <div className="mt-8 mb-10 w-full h-[600px] max-w-5xl flex items-center justify-center border border-[3px] border-[#1A535C] bg-white">
                <Graph data={dummyData} />
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
