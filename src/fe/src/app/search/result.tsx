import React from "react";
import Image from "next/image";
import Graph, { GraphData } from "../components/chart";

interface PathDetail {
  url: string;
  image: string;
  title: string;
  description: string;
}

interface Props {
  paths: string[][];
  pathsDetails: PathDetail[][];
  actualSource: string;
  actualDestination: string;
  degrees: number;
  executionTime: number;
  selectedAlgorithm: string;
  graphData: GraphData;
  articleChecked: number;
}

const PathsDisplay: React.FC<Props> = ({
  paths,
  pathsDetails,
  actualSource,
  actualDestination,
  degrees,
  executionTime,
  selectedAlgorithm,
  graphData,
  articleChecked,
}) => {
  return (
    <div className="mt-8 result flex flex-col items-center justify-center">
      <div className="w-full max-w-5xl">
        <div className="text-3xl text-center text-[#1A535C]">
          Found <strong>{paths.length} path(s)</strong> with{" "}
          <strong>{degrees} degrees</strong> of separation from{" "}
          <strong>{actualSource}</strong> to{" "}
          <strong>{actualDestination}</strong> in{" "}
          <strong>{executionTime} milliseconds</strong>. Checked a total of{" "}
          <strong>{articleChecked} articles</strong> using{" "}
          <strong>{selectedAlgorithm} Algorithm.</strong>
        </div>
      </div>
      <div className="mt-8 mb-10 w-full h-[600px] max-w-5xl flex items-center justify-center border-[3px] border-[#1A535C] bg-white">
        <Graph data={graphData} len={paths[0].length} />
      </div>
      <div className="text-[#075A5A] text-3xl font-bold mt-8">
        Individual Paths
      </div>
      <div className="flex justify-center w-full">
        <div className="flex flex-wrap justify-center gap-3">
          {pathsDetails.map((details, index) => (
            <div
              key={index}
              className="flex flex-col items-center rounded-lg border-[#1A535C] border-2 w-full max-w-[500px] h-fit mt-4 mb-8"
            >
              {details.map((detail, idx) => (
                <a
                  key={idx}
                  href={detail.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="path-card flex flex-row border-b border-[#1A535C] space-x-4 p-2 hover:bg-[#A1B6BA] w-full"
                >
                  <div className="min-w-[70px] min-h-[50px] relative rounded-lg border-[#1A535C] border">
                    <Image
                      src={detail.image || "/noimage.png"}
                      alt={detail.title}
                      layout="fill"
                      className="rounded-lg"
                    />
                  </div>
                  <div className="path-content flex flex-col text-[#075A5A]">
                    <div className="path-title items-start font-bold">
                      {detail.title}
                    </div>
                    <p className="line-clamp-2">{detail.description}</p>
                  </div>
                </a>
              ))}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default PathsDisplay;
