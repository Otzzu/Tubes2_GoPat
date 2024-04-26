import React, { useRef, useEffect } from "react";
import * as d3 from "d3";
import { SimulationNodeDatum, SimulationLinkDatum } from "d3";

// Define TypeScript types for your data
export interface Node extends SimulationNodeDatum {
  id: string;
  group: number;
  url: string;
}

export interface Link extends SimulationLinkDatum<Node> {
  source: string | Node;
  target: string | Node;
}

export interface GraphData {
  nodes: Node[];
  links: Link[];
}

const color = d3.scaleOrdinal(d3.schemeCategory10);
const createLegendData = (len: number) => {
  var data = [];
  var text = "";
  for (let i = 0; i < len; i++) {
    if (i === 0) {
      text = "Start page";
    } else if (i === len - 1) {
      text = "End page";
    } else {
      text = `${i} degree${i > 1 ? "s" : ""} away`;
    }
    data.push({
      color: color(i.toString()),
      text: text,
      group: i,
    });
  }

  return data;
};

export const parseDataForGraph = (pathsArray: string[][]): GraphData => {
  const nodes: Node[] = [];
  const links: Link[] = [];
  const nodeNameSet = new Set<string>();

  pathsArray.forEach((path) => {
    for (let i = 0; i < path.length; i++) {
      const url = path[i];
      console.log("ini url",url);
      const name = url.split("/").pop()!.replace(/_/g, " ");
      if (!nodeNameSet.has(name)) {
        nodeNameSet.add(name);
        nodes.push({ id: name, group: i ,url:url});
      }
      if (i > 0) {
        const sourceName = path[i - 1].split("/").pop()!.replace(/_/g, " ");
        links.push({
          source: sourceName,
          target: name,
        });
      }
    }
  });

  return { nodes, links };
};

export const extractTitleFromUrl = (url: string): string => {
  const title = new URL(url).pathname.split("/").pop();
  return title || "";
};

const Graph: React.FC<{ data: GraphData; len: number }> = ({ data, len }) => {
  const svgRef = useRef<SVGSVGElement>(null);
  const gRef = useRef<SVGGElement>(null);

  useEffect(() => {
    if (!svgRef.current) return;

    const width = svgRef.current.clientWidth;
    const height = svgRef.current.clientHeight;

    const svg = d3.select(svgRef.current);
    const g = d3.select(gRef.current);

    g.selectAll("*").remove();

    const zoom = d3
      .zoom<SVGSVGElement, unknown>()
      .scaleExtent([0.5, 8])
      .on("zoom", (event) => {
        g.attr("transform", event.transform);
      });

    svg.call(zoom);

    const color = d3.scaleOrdinal(d3.schemeCategory10);

    const simulation = d3
      .forceSimulation(data.nodes)
      .force(
        "link",
        d3
          .forceLink(data.links)
          .id((d) => (d as Node).id)
          .distance(100)
      ) // Increase distance between nodes
      .force("charge", d3.forceManyBody().strength(-400))
      .force("center", d3.forceCenter(width / 2, height / 2));

    const legendData = createLegendData(len);

    const legend = svg
      .append("g")
      .attr("class", "legend")
      .attr("transform", "translate(10,10)")
      .selectAll("g")
      .data(legendData)
      .enter()
      .append("g")
      .attr("class", "legend")
      .attr("transform", (d, i) => `translate(0, ${i * 20})`);

    // Draw legend colored rectangles
    legend
      .append("circle")
      .attr("cx", 10)
      .attr("cy", 9)
      .attr("r", 8)
      .attr("stroke", (d, i) => i === 0 || i === data.nodes.length - 1 ? "#1A535C" : "rgb(125, 62, 7)")
      .attr("stroke-width", 2)
      .style("fill", (d) => d.color);
      

    // Draw legend text
    legend
      .append("text")
      .attr("x", 25)
      .attr("y", 9)
      .attr("dy", ".35em")
      .text((d) => d.text)
      .style("font-family", "Poppins")
      .style("font-size", 15);

    const link = g
      .append("g")
      .attr("stroke", "#999")
      .selectAll("line")
      .data(data.links)
      .join("line")
      .attr("stroke-width", 2);

    const node = g
      .append("g")
      .selectAll<SVGCircleElement, Node>("circle")
      .data(data.nodes)
      .join("circle")
      .attr("r", (d, i) => (i === 0 || i === data.nodes.length - 1 ? 20 : 15))
      .attr("fill", (d) => color(d.group.toString()))
      .attr("stroke", (d, i) => i === 0 || i === data.nodes.length - 1 ? "#1A535C" : "rgb(125, 62, 7)")
      .attr("stroke-width", 3.5)
      .style("cursor", "pointer")
      .on("click", (event, d) => {
        window.open(d.url, "_blank");  
      });

    const labels = g
      .append("g")
      .selectAll("text")
      .data(data.nodes)
      .join("text")
      .text((d) => d.id)
      .attr("x", (d) => d.x ?? 0)
      .attr("y", (d) => d.y ?? 0)
      .style("fill", "#075A5A")
      .style("font-family", "Poppins")
      .style("font-size", 15)
      .attr("text-anchor", "middle")
      .attr("dy", "0.35em")
      .attr("dx", "4em");

    node.call(
      d3
        .drag<SVGCircleElement, Node>()
        .on("start", (event) => {
          if (!event.active) simulation.alphaTarget(0.3).restart();
          event.subject.fx = event.subject.x;
          event.subject.fy = event.subject.y;
        })
        .on("drag", (event) => {
          event.subject.fx = event.x;
          event.subject.fy = event.y;
        })
        .on("end", (event) => {
          if (!event.active) simulation.alphaTarget(0);
          event.subject.fx = null;
          event.subject.fy = null;
        })
    );

    simulation.on("tick", () => {
      link
        .attr("x1", (d) => (d.source as Node).x ?? width / 2)
        .attr("y1", (d) => (d.source as Node).y ?? height / 2)
        .attr("x2", (d) => (d.target as Node).x ?? width / 2)
        .attr("y2", (d) => (d.target as Node).y ?? height / 2);

      node
        .attr("cx", (d) => d.x ?? width / 2)
        .attr("cy", (d) => d.y ?? height / 2);

      labels
        .attr("x", (d) => d.x ?? width / 2)
        .attr("y", (d) => d.y ?? height / 2);
    });

    window.addEventListener("resize", resize);
    function resize() {
      const newWidth = svgRef.current!.clientWidth;
      const newHeight = svgRef.current!.clientHeight;

      simulation.force("center", d3.forceCenter(newWidth / 2, newHeight / 2));
      simulation.alpha(0.3).restart();
    }
    return () => {
      window.removeEventListener("resize", resize);
    };
  }, [data]);

  return (
    <svg
      ref={svgRef}
      preserveAspectRatio="xMidYMid meet"
      style={{ width: "100%", height: "100%" }}
    >
      <g ref={gRef}></g>
    </svg>
  );
};

export default Graph;
