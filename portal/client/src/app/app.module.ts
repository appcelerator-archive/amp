import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

//Services
import { UsersService } from './services/users.service';
import { DockerStacksService } from './docker-stacks/services/docker-stacks.service';
import { MenuService } from './services/menu.service';
import { AuthGuard } from './services/auth-guard.service';
import { HttpService } from './services/http.service';
import { OrganizationsService } from './organizations/services/organizations.service';
import { DockerServicesService } from './docker-stacks/services/docker-services.service';
import { DockerContainersService } from './docker-stacks/services/docker-containers.service';
import { SwarmsService } from './services/swarms.service';
import { DragService } from './services/drag.service';
import { MetricsService } from './metrics/services/metrics.service';
import { LogsService } from './logs/services/logs.service';
import { NodesService } from './nodes/services/nodes.service';
import { DashboardService } from './dashboard/services/dashboard.service';
import { ColorsService } from './dashboard/services/colors.service'

//Module
import { AppRoutingModule} from './app-routing.module';

//Directive
import { DropdownDirective } from './directives/dropdown.directive'
import { DraggableDirective } from './directives/draggable.directive'
import { DropTargetDirective } from './directives/drop-target.directive'
import { TooltipDirective } from './directives/tooltip.directive'
import { MovableDirective } from './dashboard/directives/movable.directive'

//components
import { AppComponent } from './app.component';
import { SignupComponent } from './auth/signup/signup.component';
import { SigninComponent } from './auth/signin/signin.component';
import { AuthComponent } from './auth/auth/auth.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { NodesComponent } from './nodes/nodes.component';
import { DockerStacksComponent } from './docker-stacks/docker-stacks.component';
import { PasswordComponent } from './password/password.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { PageheaderComponent } from './pageheader/pageheader.component';
import { UsersComponent } from './users/users.component';
import { AmpComponent } from './amp/amp.component';
import { SwarmsComponent } from './swarms/swarms.component';
import { LogsComponent } from './logs/logs.component';
import { MetricsComponent } from './metrics/metrics.component';
import { OrganizationsComponent } from './organizations/organizations.component';
import { OrganizationComponent } from './organizations/organization/organization.component';
import { TeamComponent } from './organizations/organization/team/team.component';
import { DockerStackDeployComponent } from './docker-stacks/docker-stack-deploy/docker-stack-deploy.component';
import { DockerServicesComponent } from './docker-stacks/docker-services/docker-services.component';
import { DockerContainersComponent } from './docker-stacks/docker-containers/docker-containers.component';
import { OrganizationCreateComponent } from './organizations/organization-create/organization-create.component';
import { TeamCreateComponent } from './organizations/organization/team/team-create/team-create.component';
import { LinesComponent } from './metrics/graph/lines/lines.component';
import { SettingsComponent } from './settings/settings/settings.component';
import { ForgotComponent } from './auth/forgot/forgot.component';
import { VerifyComponent } from './auth/verify/verify.component';
import { DGraphComponent } from './dashboard/dgraph/dgraph.component';
import { DGraphEditorComponent } from './dashboard/dgraph-editor/dgraph-editor.component';
import { DgraphAlertComponent } from './dashboard/dgraph-alert/dgraph-alert.component';

@NgModule({
  declarations: [
    //Directives
    DropdownDirective,
    DraggableDirective,
    DropTargetDirective,
    TooltipDirective,
    MovableDirective,
    //Components
    AppComponent,
    SignupComponent,
    SigninComponent,
    AuthComponent,
    DashboardComponent,
    NodesComponent,
    PasswordComponent,
    SidebarComponent,
    PageheaderComponent,
    UsersComponent,
    AmpComponent,
    SwarmsComponent,
    LogsComponent,
    MetricsComponent,
    OrganizationsComponent,
    OrganizationComponent,
    TeamComponent,
    DockerStacksComponent,
    DockerStackDeployComponent,
    DockerServicesComponent,
    DockerContainersComponent,
    OrganizationCreateComponent,
    TeamCreateComponent,
    LinesComponent,
    SettingsComponent,
    ForgotComponent,
    VerifyComponent,
    DGraphComponent,
    DGraphEditorComponent,
    DgraphAlertComponent

  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    AppRoutingModule
  ],
  providers: [
    DockerStacksService,
    DockerServicesService,
    DockerContainersService,
    UsersService,
    MenuService,
    HttpService,
    OrganizationsService,
    SwarmsService,
    DragService,
    MetricsService,
    LogsService,
    NodesService,
    DashboardService,
    ColorsService,
    AuthGuard
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
