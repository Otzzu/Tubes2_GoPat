import Image from "next/image";

export default function Navbar() {
  return (
    <nav className="bg-[#D8F0F0] text-[#075A5A] text-2xl font-semibold cursor-pointer">
      <div className="mx-auto px-4 py-6 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between">
          <div className="flex justify-end flex-1">
            <a href="./about" className="hover:text-[#FF6B6B] hover:scale-110 focus:scale-110 duration-300 transition ease-in-out delay-50">About</a>
            <a href="./GoPat.pdf" download="GoPat.pdf" className="hover:text-[#FF6B6B] hover:scale-110 focus:scale-110 duration-300 transition ease-in-out delay-50 ml-4">Docs</a>
          </div>

          <div className="flex justify-center">
            <a href="/" className="flex items-center mx-8">
              <Image src="/logo.svg" alt="logo" width={80} height={80} />
            </a>
          </div>

          <div className="flex justify-start flex-1">
            <a
              href="https://github.com/Otzzu/Tubes2_GoPat"
              className="hover:text-[#FF6B6B] hover:scale-110 focus:scale-110 duration-300 transition ease-in-out delay-50"
              target="_blank"
              rel="noopener noreferrer"
            >
              Github
            </a>
            <a href="./developers" className="hover:text-[#FF6B6B] hover:scale-110 focus:scale-110 duration-300 transition ease-in-out delay-50 ml-4">Dev</a>
          </div>
        </div>
      </div>
    </nav>
  );
}
