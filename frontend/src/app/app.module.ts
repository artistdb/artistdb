import { NgModule, APP_INITIALIZER, Injector } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { APOLLO_OPTIONS } from 'apollo-angular';
import { HttpLink } from 'apollo-angular/http';
import { ApolloClientOptions, InMemoryCache } from '@apollo/client/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ArtistComponent } from './artist/artist.component';
import { ArtistDashboardComponent } from './artist/artist-dashboard/artist-dashboard.component';
import { ArtistFormComponent } from './artist/artist-form/artist-form.component';
import { LocationComponent } from './location/location.component';
import { DynamicFormComponent } from './dynamic-form/dynamic-form.component';
import { DynamicFormFieldComponent } from './dynamic-form/dynamic-form-field/dynamic-form-field.component';
import { NavbarComponent } from './navbar/navbar.component';
import { AppConfigService } from './app-config.service';

const appInitializerFn = (config: AppConfigService) => {
  return () => {
    return config.loadAppConfig();
  };
};

const apolloInitializerFn = (httpLink: HttpLink, appConfigService: AppConfigService): ApolloClientOptions<any> => {
  return {
    link: httpLink.create({uri: appConfigService.getConfig().apiUri}),
    cache: new InMemoryCache(),
  };
}

@NgModule({
  declarations: [
    AppComponent,
    ArtistComponent,
    ArtistDashboardComponent,
    ArtistFormComponent,
    LocationComponent,
    DynamicFormComponent,
    DynamicFormFieldComponent,
    NavbarComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    ReactiveFormsModule,
    HttpClientModule
  ],
  providers: [
    AppConfigService,
    {
      provide: APOLLO_OPTIONS,
      useFactory: apolloInitializerFn,
      deps: [HttpLink, AppConfigService],
    },
    {
      provide: APP_INITIALIZER,
      useFactory: appInitializerFn,
      multi: true,
      deps: [AppConfigService]
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
