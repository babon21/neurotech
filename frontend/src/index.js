import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import Container from '@material-ui/core/Container';
import CssBaseline from '@material-ui/core/CssBaseline';
import Header from './components/Header';
import Grid from '@material-ui/core/Grid';


class App extends React.Component {

  render() {
    return (
      <React.Fragment>
        <CssBaseline />
        <Container maxWidth="lg">
          <Header/>
          {/* <main> */}
          <Grid container spacing={4}>
            Контента
          </Grid>
          {/* </main> */}
        </Container>
      </React.Fragment>
    );
  }
}

// ========================================

ReactDOM.render(
  <App />,
  document.getElementById('root')
);