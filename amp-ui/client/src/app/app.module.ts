import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

//components
import { AppComponent } from './app.component';
import { SignupComponent } from './auth/signup/signup.component';
import { SigninComponent } from './auth/signin/signin.component';
import { AuthComponent } from './auth/auth/auth.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { NodesComponent } from './nodes/nodes.component';
import { StacksComponent } from './stacks/stacks.component';
import { PasswordComponent } from './password/password.component';
import { EndpointsComponent } from './endpoints/endpoints.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { PageheaderComponent } from './pageheader/pageheader.component';
import { UsersComponent } from './users/users.component';
import { PageErrorComponent } from './page-error/page-error.component';

//Services
import { UsersService } from './services/users.service';
import { StacksService } from './services/stacks.service';
import { MenuService } from './services/menu.service';
import { AuthGuard } from './services/auth-guard.service';
//Module
import { AppRoutingModule} from './app-routing.module';
import { AmpComponent } from './amp/amp.component';

@NgModule({
  declarations: [
    AppComponent,
    SignupComponent,
    SigninComponent,
    AuthComponent,
    DashboardComponent,
    NodesComponent,
    StacksComponent,
    PasswordComponent,
    EndpointsComponent,
    SidebarComponent,
    PageheaderComponent,
    UsersComponent,
    PageErrorComponent,
    AmpComponent,
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    AppRoutingModule
  ],
  providers: [StacksService, UsersService, MenuService, AuthGuard],
  bootstrap: [AppComponent]
})
export class AppModule { }
