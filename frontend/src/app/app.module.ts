import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { APOLLO_OPTIONS } from 'apollo-angular';
import { HttpLink } from 'apollo-angular/http';
import { InMemoryCache } from '@apollo/client/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ArtistComponent } from './artist/artist.component';
import { ArtistDashboardComponent } from './artist/artist-dashboard/artist-dashboard.component';
import { ArtistFormComponent } from './artist/artist-form/artist-form.component';
import { environment } from 'src/environments/environment';
import { LocationComponent } from './location/location.component';
import { DynamicFormComponent } from './dynamic-form/dynamic-form.component';
import { DynamicFormFieldComponent } from './dynamic-form/dynamic-form-field/dynamic-form-field.component';

@NgModule({
  declarations: [
    AppComponent,
    ArtistComponent,
    ArtistDashboardComponent,
    ArtistFormComponent,
    LocationComponent,
    DynamicFormComponent,
    DynamicFormFieldComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    ReactiveFormsModule,
    HttpClientModule
  ],
  providers: [
    {
      provide: APOLLO_OPTIONS,
      useFactory: (httpLink: HttpLink) => {
        return {
          cache: new InMemoryCache(),
          link: httpLink.create({
            uri: environment.graphQLUri,
          }),
        };
      },
      deps: [HttpLink],
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
