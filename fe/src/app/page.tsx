import Search from "./search/page";
import Navbar from "./components/navbar";
export default function Home() {
  return (
    <main className="min-h-screen">
      <div className="flex flex-col min-h-screen items-center px-5 md:px-24 bg-[#68CEFF] overflow-hidden">
        <Search />
      </div>
    </main>
  );
}
