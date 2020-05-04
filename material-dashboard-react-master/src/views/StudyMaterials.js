/*eslint-disable*/
import React from "react";
import GridItem from "components/Grid/GridItem.js";
import GridContainer from "components/Grid/GridContainer.js";
import Button from "components/CustomButtons/Button.js";
import { Switch, Route, Router } from "react-router-dom";
import { hist } from '../App'
import StudyMaterialsList from "views/StudyMaterialsList";
import List from "components/DisciplineList.js";
import StudyTable from "components/StudyTable.js";

var bugs = [
  'Sign contract for "What are conference organizers afraid of?"',
  "Lines From Great Russian Literature? Or E-mails From My Boss?",
  "Flooded: One year later, assessing what was lost and what was found when a ravaging rain swept through metro Detroit",
  "Create 4 Invisible User Experiences you Never Knew About"
];


export default function StudyMaterials() {

  return (
    <Router history={hist}>
      <Switch>
        <Route exact path="/admin/study-materials/:id" component={StudyMaterialsList} />
        <GridContainer>
          <GridItem xs={12} sm={12} md={12}>
            <Button color="primary">Добавить</Button>
          </GridItem>
          <GridItem xs={12} sm={12} md={12}>
            <StudyTable
              title="Дисциплины"
              headerColor="primary"
              tables={[
                {
                  tabContent: (
                    <List
                      checkedIndexes={[0, 3]}
                      tasksIndexes={[0, 1, 2, 3]}
                      tasks={bugs}
                    />
                  )
                }
              ]}
            />
          </GridItem>
        </GridContainer>
      </Switch>
    </Router>
  );
}
