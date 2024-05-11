import RouteList from './routes/privateroutes'
import React from 'react';
import { Suspense } from 'react';


function App() {
  return (
    <Suspense fallback={<div>Loading... </div>}>
      <RouteList />
    </Suspense>
  );
}

export default App;
