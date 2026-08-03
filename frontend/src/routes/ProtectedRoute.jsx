import useAuth from "../hooks/useAuth";
import Redirect from "./Redirect";

export default function ProtectedRoute({ children }) {
  const { authenticated } = useAuth();

  if (!authenticated) {
    return <Redirect to="/login" />;
  }

  return children;
}
