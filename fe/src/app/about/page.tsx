export default function About() {
  return (
    <div className="bg-gray-200 text-gray-900 tracking-wider leading-normal justify-center flex flex-col">
      <section className="bg-cyan-800 text-white text-center p-6">
        <h1 className="text-3xl font-bold">WikiRace: BFS vs IDS</h1>
      </section>

      <section className="grid md:grid-rows-3">
        <section className="p-8 px-36 py-14 bg-gray-300">
          <p className="text-2xl font-bold mb-3">WikiRace</p>
          <p>
            WikiRace is a game where participants try to navigate from one specific Wikipedia page to another, using as few links as possible, in the shortest time. This web page used two algorithms as such Breadth-First Search (BFS) and Iterative Deepening Search (IDS) to solve WikiRace Game.
          </p>
        </section>

        <section className="p-8 px-36 py-14">
          <div className="grid md:grid-cols-2 md:divide-x-4 divide-slate-800">
            <div className="pr-8">
              <p className="text-xl font-bold mb-2">Breadth-First Search (BFS)</p>
              <p>BFS explores the graph level by level. In WikiRace, it checks all possible articles in Wikipedia from the starting article, then moves on to the links of those articles, ensuring minimal steps are taken to reach the target article.</p>
            </div>
            <div className="pl-8">
              <p className="text-xl font-bold mb-2">Iterative Deepening Search (IDS)</p>
              <p>
                IDS combines the depth-first exploration with the level-by-level exploration of BFS. It starts at a shallow depth and increases the depth limit with each iteration, making it suitable for finding the optimal path in WikiRace.
              </p>
            </div>
          </div>
        </section>

        <section className="p-8 px-36 py-[4.31rem] bg-gray-100 shadow-md">
          <h2 className="text-2xl font-bold mb-3">BFS vs IDS in WikiRace</h2>
          <p>While BFS can be faster for shallower solutions, IDS is generally better at finding the shortest path in a deep graph, as common in this WikiRace Game.</p>
        </section>
      </section>

      <footer className="bg-gray-700 text-white text-center p-4 fixed inset-x-0 bottom-0">
        <p>Play WikiRace <a href="/" className="underline">here</a>.</p>
      </footer>
    </div>
  );
}
