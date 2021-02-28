import React from 'react';

import Container from 'react-bootstrap/Container';
import './App.css';

import Search from './components/search'

const App = () => (
  <Container className="p-3">
    <Search></Search>
  </Container>
);

export default App;
