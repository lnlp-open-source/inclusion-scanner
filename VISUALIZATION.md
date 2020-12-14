Inclusion Scanner Visualization
---

## Quick Start

### Prerequisites

- Elasticsearch and Kibana running
- Project Analyzed


### Create Index Pattern

To generate a graph to analyze the data from the analyzed project, go to Kibana by going to `localhost:5601` and clicking `Visualize`. From here click `Create Index Pattern` and type in `inclusion-scanner-*`. Kibana should show that there is at least 1 index matching that pattern. Click the `next step` button and when prompted, select the `timestamp` field as "Time Field" and then click `Create Index Field`.


### Generate Graph

Go to `localhost:5601` and click `Visualize`. Click `Create new visualization`. Click `Vertical Bar` and select the data source we just created above. Click the `+ add` button to create a new bucket source. On the dropdown menu, select `X-axis`. From the dropdown for Aggregration, seelect `Terms`. From the field dropdown select `terms_abused.keyword`. Optionally add a custom label for this axis. Then click Update. 


### Save Visualization

Save the visualization to view later by clicking the `Save` button at the top left of the page. Add a title and optional description, then click `Save`. 
