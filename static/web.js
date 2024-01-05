document.addEventListener("DOMContentLoaded", () => {    
});

function myFunction(event) {
  event.preventDefault(); // Prevent form submission
  const dbfileElement = document.querySelector("#dbfile");
  const dbfile = dbfileElement.value;
  console.log(dbfile);
  
  let p1 = document.querySelector("#selected_db_file");
  p1.innerText = "Current database file: " + dbfile;
  dbfilediv.appendChild(p1);
  // Clear the input field
  dbfileElement.value = "";

  // Create a FormData object
  if (dbfile == "") {
    alert("Please enter a database file path");
    return false;
  }

  // check if id="db_tables" exists and remove it
  var db_tables = document.querySelector("#db_table");
  if (db_tables) {
    db_tables.remove();
  }

  // check if a form with id="tableForm" exists
  var tableForm = document.querySelector("#tableForm");
  if (tableForm) {
    tableForm.remove();
  }

  let formData = new FormData();
  formData.append('dbfile', dbfile);

  fetch("/viewdb", {
      method: "POST",
      body: formData
    })
    .then(response => response.json())
    .then(data => {
      console.log(data);
      var tableDiv = document.querySelector("#tableDiv");
      var tableForm = document.createElement("form");
      tableForm.setAttribute("id", "tableForm");
      tableDiv.appendChild(tableForm);
      var selectList = document.createElement("select");
      selectList.setAttribute("id", "dbtables");
      selectList.setAttribute("class", "form-select");
      selectList.setAttribute("name", "dbtables");
      tableForm.appendChild(selectList);
      tableForm.setAttribute("onchange", "queryTable(event)");

      //Create and append the options
      for (var i = 0; i < data.length; i++) {
          var option = document.createElement("option");
          option.setAttribute("value", data[i]);
          option.text = data[i];
          selectList.appendChild(option);
      }
    })
    .catch(error => console.error('Error:', error));
  };
  
function queryTable(event) {
  event.preventDefault(); // Prevent form submission
  console.log(event.value);
  const dbtableElement = document.querySelector("#dbtables");
  const dbtable = dbtableElement.value;
  console.log(dbtable);
  let tableContent = document.querySelector("#tableContent");
  tableContent.innerText = "Current table: " + dbtable;
  tableContent.innerText = dbtable;


  fetch("/table?name=" + encodeURIComponent(dbtable), {
      method: "GET",
      credentials: 'include',
    })
    .then(response => response.json())
    .then(data => {
      console.log(data);

      const refreshbutton = document.createElement("button");
      refreshbutton.setAttribute("class", "btn btn-secondary btn-sm");
      refreshbutton.setAttribute("onclick", "queryTable(event)");
      refreshbutton.innerText = "refresh";

      const tbl = document.createElement("table");
      tbl.setAttribute("id", "db_table");
      tbl.setAttribute("class", "table table-striped");
      tbl.setAttribute("style", "width:100%");
      const tblBody = document.createElement("tbody");
      const tblHead = document.createElement("thead");
      const row = document.createElement("tr");
      tbl.appendChild(tblHead);
      tblHead.appendChild(row);
      tbl.appendChild(tblBody);
      tableContent.appendChild(refreshbutton);
      tableContent.appendChild(document.createElement("br"));
      tableContent.appendChild(document.createElement("br"));
      tableContent.appendChild(tbl);

      // create table header
      //console.log("data[0] is Object", Object.keys(data[0]));
      if (data === null || data.length === 0) {
        alert("Table is empty");
        return false;
      }
      let keys = Object.keys(data[0]);
      for (var i = 0; i < keys.length; i++) {
       
        const th = document.createElement("th");
        th.innerText = "  " + keys[i] + "  ";
        row.appendChild(th);
      } 

      // create table body
      
      for (var i = 0; i < data.length; i++) {
        const row_in_body = document.createElement("tr");
        for (var j = 0; j < Object.values(data[i]).length; j++) {
            let values = Object.values(data[i]);
            console.log("value: ", values[j]);
            const cell = document.createElement("td");
            cell.innerText = values[j];
            row_in_body.appendChild(cell);
          }
          tblBody.appendChild(row_in_body); 
        }
        

    })
    .catch(error => console.error('Error:', error));};