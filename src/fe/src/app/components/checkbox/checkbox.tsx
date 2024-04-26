"use client";

import { Checkbox } from "../ui/checkbox";

interface CheckBoxProps {
    onPathChange:(path:boolean) =>void;
}

export function CheckboxWithText({onPathChange} :CheckBoxProps) {
    const handlePath =(path: boolean)=> {
        onPathChange(path);
    }
  return (
    <div className="flex flex-row items-center justify-center w-full space-x-4">
      <Checkbox id="terms1" onClick={() => handlePath(true)} />
      <label
        htmlFor="terms1"
        className="text-2xl font-semibold text-[#075A5A] leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
      >
        MultiPath
      </label>
    </div>
  );
}
