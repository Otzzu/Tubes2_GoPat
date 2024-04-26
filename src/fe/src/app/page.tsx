import Search from "./search/page";
import Navbar from "./components/navbar";
export default function Home() {
  return (
    <main className="bg-[#D8F0F0] font-poppins flex justify-center min-h-screen px-5 md:px-24 xl:px-60 overflow-hidden">
      <Search />
    </main>
  );
}
