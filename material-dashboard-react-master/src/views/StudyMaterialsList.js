/*eslint-disable*/
import React from "react";
import GridItem from "components/Grid/GridItem.js";
import GridContainer from "components/Grid/GridContainer.js";
import Button from "components/CustomButtons/Button.js";
import CustomTable from "components/CustomTable.js";
import ActionTable from "components/TableWithActions.js";

var bugs = [
  'Sign contract for "What are conference organizers afraid of?"',
  "Lines From Great Russian Literature? Or E-mails From My Boss?",
  "Flooded: One year later, assessing what was lost and what was found when a ravaging rain swept through metro Detroit",
  "Create 4 Invisible User Experiences you Never Knew About"
];


export default function StudyMaterialList() {

  return (
    <GridContainer>
      <GridItem xs={12} sm={12} md={12}>
        <Button color="primary">Добавить</Button>
      </GridItem>
      <GridItem xs={12} sm={12} md={12}>
        <CustomTable
          // title="Студенты:"
          headerColor="primary"
          tables={[
            {
              tabName: "Лекции",
              tabContent: (
                <ActionTable
                  tableHeaderColor="primary"
                  tableHead={["Name", "Student work", "Year"]}
                  // checkedIndexes={[0, 3]}
                  // tasksIndexes={[0, 1, 2, 3]}
                  tableData={[
                    ["Dakota Rice", "сетей", "Oud-Turnhout"],
                    ["Minerva Hooper", "Curaçao", "Sinaai-Waas"],
                    ["Sage Rodriguez", "Netherlands", "Baileux"],
                    ["Philip Chaney", "Korea, South", "Overland Park"],
                    ["Doris Greene", "Malawi", "Feldkirchen in Kärnten"],
                    ["Mason Porter", "Исследование методов семантического и синтаксического анализа естественного языка в системах машинного перевода", "Gloucester"]
                  ]}
                />
              )
            },
            {
              tabName: "Учебные и методические пособия",
              tabContent: (
                <ActionTable
                  tableHeaderColor="primary"
                  tableHead={["Name", "Student work", "Year"]}
                  // checkedIndexes={[0, 3]}
                  // tasksIndexes={[0, 1, 2, 3]}
                  tableData={[
                    ["Dakota Rice", "сетей", "Oud-Turnhout"],
                    ["Minerva Hooper", "Curaçao", "Sinaai-Waas"],
                    ["Sage Rodriguez", "Netherlands", "Baileux"],
                    ["Philip Chaney", "Korea, South", "Overland Park"],
                    ["Doris Greene", "Malawi", "Feldkirchen in Kärnten"],
                    ["Mason Porter", "Исследование методов семантического и синтаксического анализа естественного языка в системах машинного перевода", "Gloucester"]
                  ]}
                />
              )
            }
          ]}
        />
      </GridItem>
    </GridContainer>
  );
}
