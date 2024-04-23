import React, { useRef, useEffect } from "react";
import * as d3 from "d3";
import {
  BaseType,
  Selection,
  SimulationNodeDatum,
  SimulationLinkDatum,
} from "d3";

// Define TypeScript types for your data
interface Node extends SimulationNodeDatum {
  id: string;
  group: number;
}

interface Link extends SimulationLinkDatum<Node> {
  source: string | Node; // Can be a string ID or a Node object
  target: string | Node; // Can be a string ID or a Node object
}

interface GraphData {
  nodes: Node[];
  links: Link[];
}

const Graph: React.FC<{ data: GraphData }> = ({ data }) => {
  const svgRef = useRef<SVGSVGElement>(null);

  useEffect(() => {
    if (!svgRef.current) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll("*").remove();

    // Helper function to correctly type the node.id for d3.forceLink
    const linkForceId = (d: Node | d3.SimulationNodeDatum) => {
      if ((d as Node).id) {
        return (d as Node).id;
      }
      throw new Error("Node id not found");
    };

    const color = d3.scaleOrdinal(d3.schemeCategory10);

    const simulation = d3
      .forceSimulation(data.nodes)
      .force("link", d3.forceLink(data.links).id(linkForceId))
      .force("charge", d3.forceManyBody())
      .force(
        "center",
        d3.forceCenter(
          svgRef.current.clientWidth / 2,
          svgRef.current.clientHeight / 2
        )
      );

    const link = svg
      .append("g")
      .attr("stroke", "#999")
      .selectAll("line")
      .data(data.links)
      .join("line");

    const node = svg
      .append("g")
      .attr("stroke", "#fff")
      .attr("stroke-width", 1.5)
      .selectAll<SVGCircleElement, Node>("circle")
      .data(data.nodes)
      .join("circle")
      .attr("r", (d, i) => (i === 0 || i === data.nodes.length - 1 ? 15 : 8))
      .attr("fill", (d) => (d.group === 1 ? "red" : "blue"));

    // We need to properly type the drag behavior
    const dragBehavior: any = d3
      .drag()
      .on("start", (event: any) => {
        if (!event.active) simulation.alphaTarget(0.3).restart();
        event.subject.fx = event.subject.x;
        event.subject.fy = event.subject.y;
      })
      .on("drag", (event: any) => {
        event.subject.fx = event.x;
        event.subject.fy = event.y;
      })
      .on("end", (event: any) => {
        if (!event.active) simulation.alphaTarget(0);
        event.subject.fx = null;
        event.subject.fy = null;
      });

    node.call(dragBehavior as any);

    simulation.on("tick", () => {
      link
        .attr("x1", (d) => (d.source as Node).x ?? 0)
        .attr("y1", (d) => (d.source as Node).y ?? 0)
        .attr("x2", (d) => (d.target as Node).x ?? 0)
        .attr("y2", (d) => (d.target as Node).y ?? 0);

      node.attr("cx", (d) => d.x ?? 0).attr("cy", (d) => d.y ?? 0);
    });
  }, [data]);

  return (
    <svg
      ref={svgRef}
      viewBox="0 0 1200 600"
      preserveAspectRatio="xMidYMid meet"
      style={{ width: "100%", height: "100%" }} // Ensures the SVG is responsive
    />
  );
};

export default Graph;
