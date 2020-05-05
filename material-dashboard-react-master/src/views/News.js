/*eslint-disable*/
import React from "react";
import GridItem from "components/Grid/GridItem.js";
import GridContainer from "components/Grid/GridContainer.js";
import Button from "components/CustomButtons/Button.js";
import CustomTabs from "components/CustomTabs/CustomTabs.js";
import SimpleTable from "components/SimpleTable.js";

import { connect } from 'react-redux';
import { getNews } from '../redux/actions/newsActions'


function News(props) {

  function handleEdit(e, id) {
    alert("edit alert, id:" + id)
  }

  function handleRemove(e, id) {
    alert("remove alert, id:" + id)
  }

  function getNews() {
    if (Array.isArray(props.news.news) && props.news.news.length === 0) {
      props.getNews()
    }
  
    return {
      tabContent: (
        <SimpleTable
          handleEdit={handleEdit}
          handleRemove={handleRemove}
          tasksIndexes={[0, 1]}
          tasks={props.news.news}
        />
      )
    }
  }

  function handleAdd(e) {

  }

  return (
    <GridContainer>
      <GridItem xs={12} sm={12} md={12}>
        <Button color="primary" onClick={handleAdd}>Добавить</Button>
      </GridItem>
      <GridItem xs={12} sm={12} md={12}>
        <CustomTabs
          headerColor="primary"
          tabs={[
            getNews()
          ]}
        />
      </GridItem>
    </GridContainer>
  );
}

function mapStateToProps(state) {
  return {
    news: state.news
  };
}

const mapActionsToProps = {
  getNews
};

export default connect(
  mapStateToProps,
  mapActionsToProps
)(News);
