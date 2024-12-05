const express = require('express');
const ejs = require('ejs');
const fs = require('fs');
const path = require('path');
const app = express();
app.set('view engine', 'ejs');

// Middleware to parse URL-encoded data
app.use(express.urlencoded({ extended: true }));

app.get('/', (req, res) => {
    const categoriesParam = req.query.categories;
    if (!categoriesParam) {
        return res.status(400).send('Categories parameter is missing');
    }

    // Split the categories from the query parameter
    const targetCategories = categoriesParam.split(',').map(category => category.trim());

    // Read and parse CSV data
    const data = fs.readFileSync(path.join(__dirname, 'data.csv'), 'utf8');
    const rows = data.split('\n').map(row => row.split(';'));

    // Assume the first row is the header
    const header = rows[0];
    const records = rows.slice(1).map(row => {
        let record = {};
        header.forEach((col, index) => {
            record[col] = row[index];
        });
        return record;
    });

    // Filter records based on the target categories
    const filteredRecords = records.filter(record => targetCategories.includes(record['Category']));

    console.log('Filtered Records', filteredRecords);

    // Render the template with the filtered records
    res.render('index', { records: filteredRecords });
});

app.listen(8080, () => {
    console.log('Server running on http://localhost:8080');
});
