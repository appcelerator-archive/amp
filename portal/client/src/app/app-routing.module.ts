import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';

//components
import { AppComponent } from './app.component';
import { AmpComponent } from './amp/amp.component';
import { SignupComponent } from './auth/signup/signup.component';
import { SigninComponent } from './auth/signin/signin.component';
import { ForgotComponent } from './auth/forgot/forgot.component';
import { VerifyComponent } from './auth/verify/verify.component';
import { AuthComponent } from './auth/auth/auth.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { NodesComponent } from './nodes/nodes.component';
import { DockerStacksComponent } from './docker-stacks/docker-stacks.component';
import { PasswordComponent } from './password/password.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { PageheaderComponent } from './pageheader/pageheader.component';
import { UsersComponent } from './users/users.component';
import { SwarmsComponent } from './swarms/swarms.component';
import { LogsComponent } from './logs/logs.component';
import { MetricsComponent } from './metrics/metrics.component';
import { OrganizationsComponent } from './organizations/organizations.component';
import { OrganizationCreateComponent } from './organizations/organization-create/organization-create.component';
import { TeamComponent } from './organizations/organization/team/team.component';
import { TeamCreateComponent } from './organizations/organization/team/team-create/team-create.component';
import { OrganizationComponent } from './organizations/organization/organization.component';
import { DockerStackDeployComponent } from './docker-stacks/docker-stack-deploy/docker-stack-deploy.component';
import { DockerServicesComponent } from './docker-stacks/docker-services/docker-services.component';
import { DockerContainersComponent } from './docker-stacks/docker-containers/docker-containers.component';

//Services
import { AuthGuard } from './services/auth-guard.service';

const appRoutes : Routes = [
  { path: '', redirectTo: '/auth/signin', pathMatch: 'full'  },
  { path: 'amp', component: AmpComponent, canActivate: [AuthGuard], children: [
    { path: 'organizations', component: OrganizationsComponent, canActivate: [AuthGuard] },
    { path: 'organizations/create', component: OrganizationCreateComponent, canActivate: [AuthGuard] },
    { path: 'organizations/:orgName', component: OrganizationComponent, canActivate: [AuthGuard] },
    { path: 'organizations/:orgName/team/create', component: TeamCreateComponent, canActivate: [AuthGuard] },
    { path: 'organizations/:orgName/team/:teamName', component: TeamComponent, canActivate: [AuthGuard] },
    { path: 'dashboard', component: DashboardComponent, canActivate: [AuthGuard] },
    { path: 'stacks', component: DockerStacksComponent, canActivate: [AuthGuard] },
    { path: 'stacks/:stackName/services', component: DockerServicesComponent, canActivate: [AuthGuard] },
    { path: 'stacks/:stackName/services/:serviceId/containers', component: DockerContainersComponent, canActivate: [AuthGuard] },
    { path: 'stacks/:stackName/update', component: DockerStackDeployComponent, canActivate: [AuthGuard] },
    { path: 'stacks/deploy', component: DockerStackDeployComponent, canActivate: [AuthGuard] },
    { path: 'logs/:object/:ref', component: LogsComponent, canActivate: [AuthGuard] },
    { path: 'logs', component: LogsComponent, canActivate: [AuthGuard] },
    { path: 'metrics/:object/:type/:ref', component: MetricsComponent, canActivate: [AuthGuard] },
    { path: 'metrics', component: MetricsComponent, canActivate: [AuthGuard] },
    { path: 'nodes', component: NodesComponent, canActivate: [AuthGuard] },
    { path: 'swarms', component: SwarmsComponent, canActivate: [AuthGuard] },
    { path: 'password', component: PasswordComponent, canActivate: [AuthGuard] },
    { path: 'users/:orgName', component: UsersComponent, canActivate: [AuthGuard] },
    { path: 'users', component: UsersComponent, canActivate: [AuthGuard] },
    { path: 'signup/:id', component: SignupComponent, canActivate: [AuthGuard] },
  ]},
  { path: 'auth', component: AuthComponent, children: [
    { path: 'signin', component: SigninComponent },
    { path: 'signup', component: SignupComponent },
    { path: 'forgot', component: ForgotComponent },
    { path: 'verify/:token', component: VerifyComponent }
  ]}
  //{ path: '**', redirectTo: '/auth/signin' }
];

@NgModule({
  imports: [
    RouterModule.forRoot(appRoutes)
  ],
  exports: [RouterModule],
  declarations: []
})
export class AppRoutingModule { }
