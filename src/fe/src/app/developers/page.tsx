interface Developer {
  name: string;
  photo: string,
  githubUrl: string;
}

const developers: Developer[] = [
  {
    name: 'Mesach Harmasendro',
    photo: 'sak.jpg',
    githubUrl: 'https://github.com/Otzzu'
  },
  {
    name: 'Filbert',
    photo: 'pat.jpg',
    githubUrl: 'https://github.com/Filbert88'
  },
  {
    name: 'Ivan Hendrawan Tan',
    photo: 'van.jpg',
    githubUrl: 'https://github.com/Bodleh'
  }
];
export default function Developers() {
    return (
      <div className="h-[calc(100vh-8rem)] bg-[#7aa1a1]">
        <div className="font-bold text-4xl text-center py-8 pt-[6.5rem] px-12">
          <p>Meet Our Developers</p>
        </div>
        <div className="flex-grow flex justify-center items-center">
          <div className="grid grid-cols-3 gap-4 md:gap-12 px-12">
            {developers.map((developer, index) => (
              <div key={index} className="bg-slate-200 p-6 rounded-lg shadow-md flex flex-col items-center">
                <a href={developer.githubUrl} target="_blank" rel="noreferrer">
                  <img src={developer.photo} alt={developer.name} className="w-60 h-60 rounded-full mb-4 border-black border-2
                  transition ease-in-out delay-100 hover:-translate-y-1 hover:scale-110 duration-300" />
                </a>
                <p className="text-lg font-bold">{developer.name}</p>
              </div>
            ))}
          </div>
        </div>
        <div className="bg-[#badede] h-32 fixed bottom-0 inset-x-0">

        </div>
      </div>
    );
  }