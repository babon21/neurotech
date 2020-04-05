import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Button from '@material-ui/core/Button';
import Link from '@material-ui/core/Link';
import Container from '@material-ui/core/Container';


const useStyles = makeStyles((theme) => ({
    toolbarSecondary: {
      justifyContent: 'left',
      overflowX: 'auto',
    },
    toolbarLink: {
      padding: theme.spacing(1),
      flexShrink: 0,
    },
  }));

const sections = [
    { title: 'О нас', url: 'about-us' },
    { title: 'Новости', url: 'news' },
    { title: 'Публикации', url: 'publications' },
    { title: 'Учебно-методические материалы', url: 'educational-materials' },
    { title: 'Студенты', url: 'students' },
  ];

export default function ButtonAppBar() {
  const classes = useStyles();

  return (
    <AppBar color="default">
      <Container maxWidth="lg">
        <Toolbar component="nav" variant="dense" className={classes.toolbarSecondary} >
        {sections.map((section) => (
          <Button href={section.url}>
            <Link
              component="button"
              color="inherit"
              noWrap
              key={section.title}
              variant="body1"
              className={classes.toolbarLink}
              underline='none'
            >
              {section.title}
            </Link>
          </Button>
        ))}
        </Toolbar>
      </Container>
    </AppBar>
  );
}
