"use client";
import React, { useState, useEffect } from "react";
import { source, destination, dummyData } from "./data";
import Graph from "../components/chart";

export default function Search() {
  const [currentSourceIndex, setCurrentSourceIndex] = useState(0);
  const [currentDestinationIndex, setCurrentDestinationIndex] = useState(0);
  const [selectedAlgorithm, setSelectedAlgorithm] = useState("");
  const [showResult, setShowResult] = useState(false);
  const [sourceInput, setSourceInput] = useState("");
  const [destinationInput, setDestinationInput] = useState("");
  const [paths, setPaths] = useState([]);
  const [degrees, setDegrees] = useState(0);
  const [executionTime, setExecutionTime] = useState(0);

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

  return (
    <div className="w-full font-quicksand">
      <div className="">
        <div className="flex flex-col items-center justify-center mt-4 mb-8">
          <div className="text-4xl text-[#1A535C]">Search Algorithm</div>
          <div className="space-x-4 font-bold mt-4 text-3xl">
            <button
              className={`${
                selectedAlgorithm === "BFS" ? "text-[#1A535C]" : "text-white"
              }`}
              onClick={() => handleAlgorithmSelection("BFS")}
            >
              BFS
            </button>
            <button
              className={`${
                selectedAlgorithm === "IDS" ? "text-[#1A535C]" : "text-white"
              }`}
              onClick={() => handleAlgorithmSelection("IDS")}
            >
              IDS
            </button>
          </div>
        </div>
        <h1 className="text-center text-3xl sm:text-4xl text-[#1A535C] mt-2 font-normal">
          Find the shortest paths from
        </h1>
        <div className="flex flex-col md:flex-row justify-center items-center mt-8">
          <div className="relative w-full h-[72px] bg-[#F7FFF7] border border-[3px] border-[#1A535C]">
            <input
              className="absolute top-0 left-0 bottom-0 right-0 focus:outline-none text-center text-black text-xl sm:text-3xl"
              type="text"
              placeholder={source[currentSourceIndex]}
              value={sourceInput}
              onChange={(e) => setSourceInput(e.target.value)}
            />
          </div>
          <p className="text-2xl sm:text-3xl mx-0 md:my-0 my-4 md:mx-4 text-[#1A535C]">
            to
          </p>
          <div className="relative w-full h-[72px] bg-[#F7FFF7] border border-[3px] border-[#1A535C]">
            <input
              className="absolute top-0 left-0 bottom-0 right-0 focus:outline-none text-center text-black text-xl sm:text-3xl"
              type="text"
              placeholder={destination[currentDestinationIndex]}
              value={destinationInput}
              onChange={(e) => setDestinationInput(e.target.value)}
            />
          </div>
        </div>
      </div>
      <div className="flex flex-row justify-center items-center mt-8">
        <button
          className="bg-[#FF6B6B] w-60 h-[72px] rounded-lg border border-2 border-[#1A535C] text-2xl sm:text-3xl font-bold"
          onClick={handleSearch}
        >
          Go!
        </button>
      </div>
      {/* {showResult && (
        <div className="result flex items-center justify-center">
          <div className="w-full max-w-5xl border border-2 border-black rounded-lg mt-8">
            <div className="text-3xl text-center text-[#1A535C]">
            Found <strong>{paths.length} path(s)</strong> with{" "}
            <strong>{degrees} degrees</strong> of separation from{" "}
            <strong>{sourceInput}</strong> to{" "}
            <strong>{destinationInput}</strong>
            <strong>{executionTime} seconds</strong> using{" "}
            <strong>{selectedAlgorithm} Algorithm.</strong>
          </div>
          <div className="w-full h-full max-w-5xl flex items-center justify-center border border-black bg-white">
            <Graph data={dummyData} />
          </div>
        </div>
      )} */}
      <div className="result flex flex-col items-center justify-center">
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
    </div>
  );
}
