import useAuth from "../hooks/useAuth";
import Redirect from "./Redirect";

export default function RoleRoute({ roles, children }) {
  const { user } = useAuth();

  if (!user || !roles.includes(user.role)) {
    return <Redirect to="/forbidden" />;
  }

  return children;
}
