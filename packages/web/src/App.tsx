import { QueryClientProvider, QueryClient } from "@tanstack/react-query";
import { useDisconnect, useProfile } from "./stores/profile";

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
  const disconnect = useDisconnect();
  return <button onClick={disconnect}>disconnect</button>;
}

function Profile() {
  const profile = useProfile();
  return <pre>{JSON.stringify(profile.data, null, 2)}</pre>;
}

export default App;
