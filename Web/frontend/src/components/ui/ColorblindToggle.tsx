"use client";

import { useEffect, useState } from "react";

import { Button } from "@/components/ui/Button";

const STORAGE_KEY = "vision-mode";

export function ColorblindToggle() {
  const [isTritanopia, setIsTritanopia] = useState(false);

  useEffect(() => {
    const saved = window.localStorage.getItem(STORAGE_KEY);
    if (saved === "tritanopia" || saved === "deuteranopia") {
      setIsTritanopia(true);
    }
  }, []);

  useEffect(() => {
    const root = document.documentElement;

    if (isTritanopia) {
      root.dataset.vision = "tritanopia";
      window.localStorage.setItem(STORAGE_KEY, "tritanopia");
    } else {
      root.removeAttribute("data-vision");
      window.localStorage.removeItem(STORAGE_KEY);
    }
  }, [isTritanopia]);

  return (
    <div className="fixed right-4 bottom-4 z-50">
      <Button
        variant="secondary"
        type="button"
        aria-pressed={isTritanopia}
        onClick={() => setIsTritanopia((current) => !current)}
        className="shadow-md"
      >
        {isTritanopia ? "Tritanopia enabled" : "Tritanopia"}
      </Button>
    </div>
  );
}
