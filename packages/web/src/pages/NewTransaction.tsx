import { createRoute, Link } from "@tanstack/react-router";
import { useDisconnect } from "@/stores/profile";
import { root } from "./Root";

function NewTransaction() {
  const disconnect = useDisconnect();

  return (
    <>
      <p>new</p>
      <Link to="/">go to index</Link>
      <button className="block" onClick={disconnect}>
        disconnect
      </button>
    </>
  );
}

export const newTransaction = createRoute({
  getParentRoute: () => root,
  path: "/new",
  component: NewTransaction,
});
