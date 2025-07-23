function App() {
  function disconnect() {
    fetch("http://localhost:8888/oauth2/revoke", {
      method: "post",
      credentials: "include",
    });
  }

  return (
    <>
      <h1>Blingpot</h1>
      <a className="block" href="http://localhost:8888/oauth2/authenticate">
        authenticate
      </a>
      <button onClick={disconnect}>disconnect</button>
    </>
  );
}

export default App;
