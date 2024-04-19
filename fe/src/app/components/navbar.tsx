import Link from "next/link";
import Image from "next/image";
export default function Navbar() {
  return (
    <nav className="bg-[#68CEFF] text-black">
      <div className="max-w-7xl mx-auto px-4 py-6 pt-10 sm:px-6 lg:px-8">
        <div className="flex justify-center">
          <div className="flex">
            <div className="flex-shrink-0 flex flex-col items-center">
              <a href="/">
                <div className="relative sm:w-[450px] sm:h-[100px] w-[250px] h-[70px]">
                  <Image
                    src="/logo.png"
                    alt="logo"
                    layout="fill"
                  />
                </div>
              </a>
              <p className="font-quicksand text-2xl font-bold text-[#1A535C] mt-2">
                by GoPat
              </p>
            </div>
          </div>
        </div>
        <div className="flex flex-row justify-center mt-4 space-x-2 sm:space-x-4 items-center text-lg sm:text-xl font-quicksand font-bold text-[#1A535C] cursor-pointer">
          <a className="hover:text-[#FF6B6B] hover:scale-110 focus:scale-110 duration-300 transition ease-in-out delay-50">About</a>
          <div className="h-8 bg-black w-0.5"></div>
          <a className="hover:text-[#FF6B6B] hover:scale-110 focus:scale-110 duration-300 transition ease-in-out delay-50">Developers</a>
          <div className="h-8 bg-black w-0.5"></div>
          <a className="hover:text-[#FF6B6B] hover:scale-110 focus:scale-110 duration-300 transition ease-in-out delay-50" href="https://github.com/Otzzu/Tubes2_GoPat" target="blank">Github</a>
        </div>
      </div>
    </nav>
  );
}
