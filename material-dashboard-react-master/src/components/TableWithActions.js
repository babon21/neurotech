import React from "react";
import PropTypes from "prop-types";
import classnames from "classnames";
// @material-ui/core components
import { makeStyles } from "@material-ui/core/styles";
import Checkbox from "@material-ui/core/Checkbox";
import Tooltip from "@material-ui/core/Tooltip";
import IconButton from "@material-ui/core/IconButton";
import Table from "@material-ui/core/Table";
import TableRow from "@material-ui/core/TableRow";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
// @material-ui/icons
import Edit from "@material-ui/icons/Edit";
import Close from "@material-ui/icons/Close";
import Check from "@material-ui/icons/Check";
// core components
import styles from "assets/jss/material-dashboard-react/components/tasksStyle.js";
// import styles from "assets/jss/material-dashboard-react/components/tableStyle.js";

import TableHead from "@material-ui/core/TableHead";
// import TableHead from "@material-ui/core/TableHead";
// import TableRow from "@material-ui/core/TableRow";
// import TableRow from "@material-ui/core/TableRow";


const useStyles = makeStyles(styles);

export default function Tasks(props) {
  const classes = useStyles();
  // const [checked, setChecked] = React.useState([...props.checkedIndexes]);
  // const handleToggle = value => {
  //   const currentIndex = checked.indexOf(value);
  //   const newChecked = [...checked];
  //   if (currentIndex === -1) {
  //     newChecked.push(value);
  //   } else {
  //     newChecked.splice(currentIndex, 1);
  //   }
  //   setChecked(newChecked);
  // };
  const { tableHead, tableHeaderColor, tasksIndexes, tableData, rtlActive } = props;
  // const tableCellClasses = classnames(classes.tableCell, {
  //   [classes.tableCellRTL]: rtlActive
  // });
  return (
    <Table className={classes.table}>
      {tableHead !== undefined ? (
        <TableHead className={classes[tableHeaderColor + "TableHeader"]}>
          <TableRow className={classes.tableHeadRow}>
            {tableHead.map((prop, key) => {
              return (
                <TableCell
                  className={classes.tableCell + " " + classes.tableHeadCell}
                  key={key}
                >
                  {prop}
                </TableCell>
              );
            })}
          </TableRow>
        </TableHead>
      ) : null}
      <TableBody>
        {tableData.map((prop, key) => {
          return (
            <TableRow key={key} className={classes.tableRow}>
              {prop.map((prop, key) => {
                return (
                  <TableCell className={classes.tableCell} key={key}>
                    {prop}
                  </TableCell>
                );
              })}
              <TableCell className={classes.tableActions}>
                <Tooltip
                  id="tooltip-top"
                  title="Edit Task"
                  placement="top"
                  classes={{ tooltip: classes.tooltip }}
                >
                  <IconButton
                    aria-label="Edit"
                    className={classes.tableActionButton}
                  >
                    <Edit
                      className={
                        classes.tableActionButtonIcon + " " + classes.edit
                      }
                    />
                  </IconButton>
                </Tooltip>
                <Tooltip
                  id="tooltip-top-start"
                  title="Remove"
                  placement="top"
                  classes={{ tooltip: classes.tooltip }}
                >
                  <IconButton
                    aria-label="Close"
                    className={classes.tableActionButton}
                  >
                    <Close
                      className={
                        classes.tableActionButtonIcon + " " + classes.close
                      }
                    />
                  </IconButton>
                </Tooltip>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}

Tasks.propTypes = {
  tasksIndexes: PropTypes.arrayOf(PropTypes.number),
  tasks: PropTypes.arrayOf(PropTypes.node),
  rtlActive: PropTypes.bool,
  checkedIndexes: PropTypes.array
};
