import { useEffect, useState } from "react";

export default function usePathname() {
  const [pathname, setPathname] = useState(window.location.pathname);

  useEffect(() => {
    function handleLocationChange() {
      setPathname(window.location.pathname);
    }

    window.addEventListener("popstate", handleLocationChange);
    return () => window.removeEventListener("popstate", handleLocationChange);
  }, []);

  return pathname;
}
