import {
  Tabs,
  TabsList,
  TabsTrigger,
} from "../ui/tabs"

interface TabsDemoProps {
    onAlgorithmChange: (algorithm: string) => void;
}

export function TabsDemo({ onAlgorithmChange } : TabsDemoProps) {
    const handleTabClick = (algorithm : string) => {
      onAlgorithmChange(algorithm);
    };
  
    return (
      <Tabs defaultValue="account" className="w-[300px] sm:w-[705px]">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="account" className="text-[#A1B6BA] text-3xl font-semibold" onClick={() => handleTabClick('BFS')}>BFS</TabsTrigger>
          <TabsTrigger value="password" className="text-[#A1B6BA] text-3xl font-semibold" onClick={() => handleTabClick('IDS')}>IDS</TabsTrigger>
        </TabsList>
      </Tabs>
    )
  }
  
