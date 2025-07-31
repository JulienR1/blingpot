import { createRoute, Link } from "@tanstack/react-router";
import { root } from ".";
import { useDisconnect } from "@/stores/profile";

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
