/*eslint-disable*/
import React from "react";
import GridItem from "components/Grid/GridItem.js";
import GridContainer from "components/Grid/GridContainer.js";
import Button from "components/CustomButtons/Button.js";
import CustomTabs from "components/CustomTabs/CustomTabs.js";
import Tasks from "components/SimpleTable.js";

var bugs = [
  'Sign contract for "What are conference organizers afraid of?"',
  "Lines From Great Russian Literature? Or E-mails From My Boss?",
  "Flooded: One year later, assessing what was lost and what was found when a ravaging rain swept through metro Detroit",
  "Create 4 Invisible User Experiences you Never Knew About"
];

export default function News() {

  return (
    <GridContainer>
      <GridItem xs={12} sm={12} md={12}>
        <Button color="primary">Добавить</Button>
      </GridItem>
      <GridItem xs={12} sm={12} md={12}>
        <CustomTabs
          headerColor="primary"
          tabs={[
            {
              tabContent: (
                <Tasks
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
  );
}
