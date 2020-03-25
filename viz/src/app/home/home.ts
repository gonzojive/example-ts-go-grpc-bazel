import {Component, NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import { CommonModule } from '@angular/common';

import { PlotlyViaWindowModule } from 'angular-plotly.js';
import { Plotly } from 'angular-plotly.js/src/app/shared/plotly.interface';
//import { Plotly } from 'plotly.js-mapbox-dist';
//import { Plotly } from 'plotly.js-basic-dist';
//import * as fepb from 'examples_angular/httpserver/frontendpb/frontend_pb';
//import * as feservice from 'examples_angular/httpserver/frontendpb/frontend_pb_service';
import * as feservice from 'examples_angular/httpserver/frontendpb/frontend_grpc_web_pb';
import { GetGraphRequest, GetGraphResponse } from 'examples_angular/httpserver/frontendpb/frontend_pb';
//import * as xx from 'google-protobuf';

//PlotlyModule.plotlyjs = Plotly;


const GRAPH2 = {
  data: [
      { x: [1, 2, 3], y: [2, 1, 1], type: 'scatter', mode: 'lines+points', marker: {color: 'purple'} },
      { x: [1, 2, 3], y: [2, 5, 3], type: 'bar' },
  ],
  layout: {autosize: true, title: 'A Fancy Plot'}
};

const GRAPH_SUCCESS = {
  data: [
      { x: [1, 2, 3], y: [2, 1, 1], type: 'scatter', mode: 'lines+points', marker: {color: 'green'} },
      { x: [1, 2, 3], y: [2, 5, 3], type: 'bar' },
  ],
  layout: {autosize: true, title: 'A Fancy Plot'}
};

@Component({
  selector: 'home',
  templateUrl: './home.html',
})
export class Home {
  public graph = {
    data: [
        { x: [1, 2, 3], y: [2, 6, 3], type: 'scatter', mode: 'lines+points', marker: {color: 'red'} },
        { x: [1, 2, 3], y: [2, 5, 3], type: 'bar' },
    ],
    layout: {autosize: true, title: 'A Fancy Plot'}
  };

  ngOnInit() {
    console.log("client class: %o", feservice.FrontendServicePromiseClient);
    const serviceUrl = '/';
    const client: feservice.FrontendServicePromiseClient = new feservice.FrontendServicePromiseClient(serviceUrl);
    console.log("client object: %o", client);
    const request = new GetGraphRequest();
    request.setSeed("abc");
    client.getGraph(request).then((response: GetGraphResponse) => {
      this.graph = GRAPH_SUCCESS;
    }, () => {
      this.graph = GRAPH2;
    });
    //console.log("plotly: %o", Plotly);
  }
};

@NgModule({
  declarations: [Home],
  imports: [
    CommonModule,
    PlotlyViaWindowModule,
    //PlotlyModule,
    RouterModule.forChild([{path: '', component: Home}]),
  ],
})


export class HomeModule {

}