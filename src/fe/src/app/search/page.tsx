"use client";
import React, { useState, useEffect } from "react";
import Image from "next/image";
import { source, destination } from "./data";
import { CheckboxWithText } from "../components/checkbox/checkbox";
import {
  parseDataForGraph,
  GraphData,
  extractTitleFromUrl,
} from "../components/chart";
import { TabsDemo } from "../components/tabs/tabs";
import PathsDisplay from "./result";

interface SearchResult {
  title: string;
  url: string;
}

interface SearchResultDetails {
  title: string;
  url: string;
  description: string;
  image: string;
}

export default function Search() {
  const [isMultiPath, setIsMultiPath] = useState(false);
  const [isLoading, setisLoading] = useState(false);
  const [isRotating, setIsRotating] = useState(false);
  const [selectedAlgorithm, setSelectedAlgorithm] = useState("BFS");
  const [showResult, setShowResult] = useState(false);

  const [currentSourceIndex, setCurrentSourceIndex] = useState(0);
  const [currentDestinationIndex, setCurrentDestinationIndex] = useState(0);

  const [sourceInput, setSourceInput] = useState("");
  const [sourceResults, setSourceResults] = useState<SearchResult[]>([]);
  const [sourceFocused, setSourceFocused] = useState(false);
  const [sourceURL, setSourceURL] = useState("");
  const [actualSource, setActualSource] = useState("");
  const [actualDestination, setActualDestination] = useState("");

  const [destinationInput, setDestinationInput] = useState("");
  const [destinationResults, setDestinationResults] = useState<SearchResult[]>(
    []
  );
  const [destinationFocused, setDestinationFocused] = useState(false);
  const [destinationURL, setDestinationURL] = useState("");

  const [graphData, setGraphData] = useState<GraphData>({
    nodes: [],
    links: [],
  });
  const [paths, setPaths] = useState<string[][]>([]);
  const [degrees, setDegrees] = useState(0);
  const [executionTime, setExecutionTime] = useState(0);
  const [pathsDetails, setPathsDetails] = useState<SearchResultDetails[][]>([]);
  const [articleChecked, setarticleChecked] = useState(0);

  const fetchResults = async (
    input: string,
    setResult: React.Dispatch<React.SetStateAction<SearchResult[]>>
  ) => {
    if (input) {
      try {
        const response = await fetch(
          `https://en.wikipedia.org/w/api.php?action=opensearch&search=${input}&format=json&origin=*`
        );
        const data = await response.json();
        const formattedResults = data[1].map(
          (title: string, index: number) => ({
            title: title,
            url: data[3][index],
          })
        );
        setResult(formattedResults);
      } catch (error) {
        console.error("Failed to fetch data:", error);
      }
    } else {
      setResult([]);
    }
  };

  const fetchAdditionalDetails = async (paths: string[][]) => {
    console.log(paths);
    const details: SearchResultDetails[][] = [];

    for (let i = 0; i < paths.length; i++) {
      const detail: SearchResultDetails[] = await Promise.all(
        paths[i].map(async (path) => {
          const title = extractTitleFromUrl(path);
          const response = await fetch(
            `https://en.wikipedia.org/api/rest_v1/page/summary/${title}`
          );
          const data = await response.json();
          return {
            title: data.title,
            url: path,
            description: data.extract,
            image: data.thumbnail ? data.thumbnail.source : "",
          };
        })
      );

      details.push(detail);
    }
    setPathsDetails(details);
    console.log("ini detail", details);
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      fetchResults(sourceInput, setSourceResults);
    }, 300);
    console.log(sourceInput);
    return () => clearTimeout(timer);
  }, [sourceInput]);

  useEffect(() => {
    const timer = setTimeout(() => {
      fetchResults(destinationInput, setDestinationResults);
    }, 300);
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

  const togglePathType = () => {
    setIsMultiPath((prevIsMultiPath) => !prevIsMultiPath);
  };

  const handleSearchMulti = async () => {
    console.log("multi");
    setShowResult(false);
    if (!selectedAlgorithm || !sourceInput || !destinationInput) {
      alert(
        "Please select an algorithm and enter both a source and a destination."
      );
      return;
    }
    setActualSource(sourceInput);
    setActualDestination(destinationInput);
    setisLoading(true);

    try {
      const response = await fetch(`http://localhost:8080/search/BFS/multi`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ start: sourceURL, goal: destinationURL }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      console.log("Data received:", data);
      const parsedGraphData = parseDataForGraph(data.paths);
      await fetchAdditionalDetails(data.paths);

      setGraphData(parsedGraphData);
      setShowResult(true);
      setPaths(data.paths);
      setDegrees(data.paths[0].length - 1);
      setExecutionTime(data.executionTime);
      setarticleChecked(data.countChecked);
    } catch (error) {
      console.error("There was an error!", error);
    } finally {
      setisLoading(false);
    }
  };

  useEffect(() => {
    const data = isMultiPath;
    console.log(data);
  }, [isMultiPath]);

  const handleSearchSingle = async () => {
    setShowResult(false);
    if (!selectedAlgorithm || !sourceInput || !destinationInput) {
      alert(
        "Please select an algorithm and enter both a source and a destination."
      );
      return;
    }
    setActualSource(sourceInput);
    setActualDestination(destinationInput);
    setisLoading(true);

    try {
      const response = await fetch(
        `http://localhost:8080/search/${selectedAlgorithm}2`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ start: sourceURL, goal: destinationURL }),
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      console.log("Data received:", data);
      const parsedGraphData = parseDataForGraph(data.paths);
      setGraphData(parsedGraphData);

      // Update state with the response data
      setShowResult(true);
      setPaths(data.paths);
      await fetchAdditionalDetails(data.paths);
      setDegrees(data.paths[0].length - 1);
      setExecutionTime(data.executionTime);
      setarticleChecked(data.countChecked);
    } catch (error) {
      console.error("There was an error!", error);
    } finally {
      setisLoading(false);
    }
  };

  const swapInputs = () => {
    setIsRotating(true);
    setTimeout(() => {
      setIsRotating(false);
    }, 500);

    setSourceInput(destinationInput);
    setDestinationInput(sourceInput);
    setSourceURL(destinationURL);
    setDestinationURL(sourceURL);
  };

  return (
    <div className="w-full flex flex-col justify-center items-center">
      <div className="max-w-3xl w-full h-full ">
        <div className="min-w-full my-8">
          <div className="bg-white p-8 rounded-xl">
            <div className="flex flex-col items-start mb-8">
              <div className="text-3xl sm:text-4xl text-[#075A5A] font-bold">
                Search Algorithm
              </div>
              <div className="mt-4 mb-4 flex flex-col w-full">
                <TabsDemo onAlgorithmChange={handleAlgorithmSelection} />
              </div>
              <div className="h-[1px] mt-4 w-full bg-[#A1B6BA]"></div>
              {selectedAlgorithm === "BFS" && (
                <div className="mt-8 w-full">
                  <CheckboxWithText onPathChange={togglePathType} />
                </div>
              )}
            </div>
            <div className="flex flex-col bg-white ">
              <div className="bg-[#E4F8F8] relative w-full h-[120px] py-4 rounded-lg px-8">
                <div className="text-[#3F6870] text-3xl flex flex-row justify-items-start">
                  From
                </div>
                <input
                  className="absolute bg-[#E4F8F8] bottom-4 left-8 right-0 focus:outline-none text-[#075A5A] font-semibold text-xl sm:text-3xl"
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
                          e.preventDefault();
                        }}
                        onClick={() => {
                          setSourceInput(result.title);
                          setSourceURL(result.url);
                          setSourceResults([]);
                          setSourceFocused(false);
                        }}
                      >
                        {result.title}
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
                  className="absolute bg-[#E4F8F8] bottom-4 left-8 right-0 focus:outline-none text-[#075A5A] font-semibold text-xl sm:text-3xl"
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
                          setDestinationInput(result.title);
                          setDestinationURL(result.url);
                          setDestinationResults([]);
                          setDestinationFocused(false);
                        }}
                      >
                        {result.title}
                      </li>
                    ))}
                  </ul>
                )}
              </div>
            </div>
            <div className="flex flex-row justify-center items-center mt-8">
              <button
                className="bg-[#075A5A] text-white w-full h-[72px] rounded-lg  border-2 border-[#1A535C] text-2xl sm:text-3xl font-bold"
                onClick={isMultiPath ? handleSearchMulti : handleSearchSingle}
              >
                Go!
              </button>
            </div>
            {isLoading && (
              <div className="fixed inset-0 bg-black bg-opacity-80 flex justify-center items-center z-100 no-doc-scroll">
                <div className="bg-white rounded-lg max-w-xs w-full border-2 border-[#A1B6BA]">
                  <div className="pt-8 text-[#075A5A] text-xl flex flex-col justify-center items-center text-center space-y-8">
                    <div className="loader flex justify-center"></div>
                    <div className="font-bold py-2">Please kindly wait</div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
      {showResult && (
        <PathsDisplay
          paths={paths}
          pathsDetails={pathsDetails}
          actualSource={actualSource}
          actualDestination={actualDestination}
          degrees={degrees}
          executionTime={executionTime}
          selectedAlgorithm={selectedAlgorithm}
          graphData={graphData}
          articleChecked={articleChecked}
        />
      )}
    </div>
  );
}
