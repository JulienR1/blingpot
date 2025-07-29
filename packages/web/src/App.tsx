import {
  QueryClientProvider,
  QueryClient,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { useCallback } from "react";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <h1>Blingpot</h1>
      <a className="block" href="/oauth2/authenticate">
        authenticate
      </a>
      <Disconnect />
      <Profile />
    </QueryClientProvider>
  );
}

function Disconnect() {
  const q = useQueryClient();

  const disconnect = useCallback(async () => {
    await fetch("/oauth2/revoke", { method: "post", credentials: "include" });
    q.invalidateQueries({ queryKey: ["user-profile"] });
  }, [q]);

  return <button onClick={disconnect}>disconnect</button>;
}

function Profile() {
  const profile = useQuery({
    queryKey: ["user-profile"],
    queryFn: () =>
      fetch("/profiles/me", { credentials: "include" })
        .then((res) => res.json())
        .catch(() => null),
  });

  return <pre>{JSON.stringify(profile.data, null, 2)}</pre>;
}

export default App;
