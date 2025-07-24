function App() {
  function disconnect() {
    fetch(`${import.meta.env.VITE_SERVER_URL}/oauth2/revoke`, {
      method: "post",
      credentials: "include",
    });
  }

  return (
    <>
      <h1>Blingpot</h1>
      <a
        className="block"
        href={`${import.meta.env.VITE_SERVER_URL}/oauth2/authenticate`}
      >
        authenticate
      </a>
      <button onClick={disconnect}>disconnect</button>
    </>
  );
}

export default App;
