import { useEffect } from "react";
import { navigate } from "../utils/navigation";

export default function Redirect({ to, replace = true }) {
  useEffect(() => {
    navigate(to, { replace });
  }, [to, replace]);

  return null;
}
