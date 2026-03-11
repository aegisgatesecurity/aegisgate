const fs = require('fs');
const path = 'C:/Users/Administrator/Desktop/Testing/aegisgate/.github/workflows/ml-pipeline.yml';
const content = fs.readFileSync(path, 'utf8');
console.log(content);
