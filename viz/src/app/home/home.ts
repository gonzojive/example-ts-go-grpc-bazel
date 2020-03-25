import {Component, NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import { CommonModule } from '@angular/common';

import { PlotlyViaWindowModule } from 'angular-plotly.js';
import { Plotly } from 'angular-plotly.js/src/app/shared/plotly.interface';
import * as feservice from 'examples_angular/httpserver/frontendpb/frontend_grpc_web_pb';
import { GetGraphRequest, GetGraphResponse } from 'examples_angular/httpserver/frontendpb/frontend_pb';

@Component({
  selector: 'home',
  templateUrl: './home.html',
})
export class Home {
  public graph: Plotly.Figure = {
    data: [
        { x: [1, 2, 3], y: [2, 6, 3], type: 'scatter', mode: 'lines+points', marker: {color: 'red'} },
        { x: [1, 2, 3], y: [2, 5, 3], type: 'bar' },
    ],
    layout: {autosize: true, title: 'Initial plot'},
    frames: null
  };

  ngOnInit() {
    const serviceUrl = window.document.location.protocol + '//' + window.document.location.host;
    const client: feservice.FrontendServicePromiseClient = new feservice.FrontendServicePromiseClient(serviceUrl);
    const request = new GetGraphRequest();
    request.setScale(Math.random());
    client.getGraph(request).then((response: GetGraphResponse) => {
      console.log("got response from server: %o", request);
      this.graph = responseToPlotlyGraph(response);
    }, (err) => {
      console.error("error requesting graph: %o", err);
    });
  }
};

function responseToPlotlyGraph(request: GetGraphResponse): Plotly.Figure {
  return {
    data: request.getTracesList().map((trace): Plotly.Data => {
      return {
        x: trace.getPointsList().map((pt) => pt.getX()),
        y: trace.getPointsList().map((pt) => pt.getY()),
        mode: trace.getPlotlyMode() || 'lines+points',
        type: trace.getPlotlyType() || 'scatter'
      };
    }),
    layout: {autosize: true, title: 'Graph from gRPC'},
    frames: null
  };
}

@NgModule({
  declarations: [Home],
  imports: [
    CommonModule,
    PlotlyViaWindowModule,
    RouterModule.forChild([{path: '', component: Home}]),
  ],
})


export class HomeModule {

}