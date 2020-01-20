require("expose-loader?$!expose-loader?jQuery!jquery");
require("papaparse/papaparse.min.js");

// PapaParse parser Config
var parserConfig = {
        header: true,
        complete: function(res) {
                createData(res.data[0]);
        }
}

// add event listener/ handler
var _btn = document.getElementById('import_btn');
_btn.addEventListener('click', getCSVFile);

// getCSVFile gets the uploaded CSV file which the users has added
function getCSVFile() {
        console.log("Getting CSV file");

        var file = document.getElementById('input-file').files[0];
        parseCSV(file);

        console.log("[----] CSV parsed!");
}

function parseCSV(file) {
        Papa.parse(url, parserConfig);
}

function createData(data) {
        console.log("Data: "+ data);

}


