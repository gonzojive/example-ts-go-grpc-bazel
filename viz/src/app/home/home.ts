import {Component, NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import { CommonModule } from '@angular/common';

import { PlotlyViaWindowModule } from 'angular-plotly.js';
import { Plotly } from 'angular-plotly.js/src/app/shared/plotly.interface';
import * as fepb from 'examples_angular/httpserver/frontendpb/frontend_pb';
//import * as feservice from 'examples_angular/httpserver/frontendpb/frontend_pb_service';
import * as feservice from 'examples_angular/httpserver/frontendpb/frontend_pb_service';
import * as xx from 'google-protobuf';


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
    //console.log("initialized home comonent w/protobuf %o", xx);
    const service: feservice.FrontendService = new feservice.FrontendService();
  }
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