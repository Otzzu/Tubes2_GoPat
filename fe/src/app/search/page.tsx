"use client";
import React, { useState, useEffect } from "react";

export default function Search() {
  const source = [
    "Donald Trump",
    "Beale Cipher",
    "Einstein Theory",
    "Python Basics",
    "Mars Rover",
    "Artificial Intelligence",
    "Blockchain Technology",
    "Quantum Computing",
    "Renaissance Art",
    "Silk Road",
    "Climate Change",
    "Electric Cars",
    "Solar Power",
    "Music Theory",
    "Ancient Egypt",
    "Deep Learning",
    "Startup Culture",
    "Mount Everest",
    "Hollywood Cinema",
    "Digital Marketing",
    "Viking History",
    "Robotics Innovation",
  ];

  const destination = [
    "Machine Learning",
    "Galaxy Formation",
    "Classical Music",
    "Sustainable Farming",
    "Cyber Security",
    "Literary Critique",
    "Space Tourism",
    "Game Development",
    "Mental Health",
    "Public Speaking",
    "Financial Freedom",
    "Quantum Mechanics",
    "Ancient Rome",
    "Virtual Reality",
    "Coral Reefs",
    "Quantum Entanglement",
    "Data Privacy",
    "Nutrition Facts",
    "Clean Energy",
    "Urban Planning",
  ];

  const [currentSourceIndex, setCurrentSourceIndex] = useState(0);
  const [currentDestinationIndex, setCurrentDestinationIndex] = useState(0);

  useEffect(() => {
    const updatePlaceholders = () => {
      setCurrentSourceIndex(
        (prevIndex) => (prevIndex + 1) % source.length
      );

      setCurrentDestinationIndex(
        (prevIndex) => (prevIndex + 1) % destination.length
      );
    };

    const intervalId = setInterval(updatePlaceholders, 2000);

    return () => clearInterval(intervalId);
  }, []);

  return (
    <div className="w-full font-quicksand">
      <div className="">
        <h1 className="text-center text-3xl sm:text-4xl text-[#1A535C] mt-2 font-normal">
          Find the shortest paths from
        </h1>
        <div className="flex flex-col md:flex-row justify-center items-center mt-8">
          <div className="relative w-full max-w-lg h-[72px] bg-[#F7FFF7] border border-[3px] border-[#1A535C]">
            <input
              className="absolute top-0 left-0 bottom-0 right-0 focus:outline-none text-center text-black text-xl sm:text-3xl"
              type="text"
              placeholder={source[currentSourceIndex]}
            ></input>
          </div>
          <p className="text-2xl sm:text-3xl mx-0 md:my-0 my-4 md:mx-4 text-[#1A535C]">
            to
          </p>
          <div className="relative w-full max-w-lg h-[72px] bg-[#F7FFF7] border border-[3px] border-[#1A535C]">
            <input
              className="absolute top-0 left-0 bottom-0 right-0 focus:outline-none text-center text-black text-xl sm:text-3xl"
              type="text"
              placeholder={destination[currentSourceIndex]}
            ></input>
          </div>
        </div>
      </div>
      <div className="flex flex-row justify-center items-center mt-8">
        <button className="bg-[#FF6B6B] w-60 h-[72px] rounded-lg border border-2 border-[#1A535C] text-2xl sm:text-3xl font-bold">
          Go!
        </button>
      </div>
    </div>
  );
}
