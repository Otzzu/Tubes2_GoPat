import Search from "./search/page";
import Navbar from "./components/navbar";
export default function Home() {
  return (
    <main className="min-h-screen">
      <div className="flex flex-col min-h-screen items-center justify-center px-5 md:px-24 xl:px-60 bg-[#D8F0F0] overflow-hidden">
        <Search />
      </div>
    </main>
  );
}
