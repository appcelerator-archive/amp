import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Routes, RouterModule } from '@angular/router';

//components
import { AppComponent } from './app.component';
import { AmpComponent } from './amp/amp.component';
import { SignupComponent } from './auth/signup/signup.component';
import { SigninComponent } from './auth/signin/signin.component';
import { AuthComponent } from './auth/auth/auth.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { NodesComponent } from './nodes/nodes.component';
import { StacksComponent } from './stacks/stacks.component';
import { PasswordComponent } from './password/password.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { PageheaderComponent } from './pageheader/pageheader.component';
import { UsersComponent } from './users/users.component';

//Services
import { AuthGuard } from './services/auth-guard.service';

const appRoutes : Routes = [
  { path: '', redirectTo: '/auth/signin', pathMatch: 'full'  },
  { path: 'amp', component: AmpComponent, canActivate: [AuthGuard], children: [
    { path: 'dashboard', component: DashboardComponent, canActivate: [AuthGuard] },
    { path: 'nodes', component: NodesComponent, canActivate: [AuthGuard] },
    { path: 'stacks', component: StacksComponent, children: [
      { path: ':name', component: StacksComponent, canActivate: [AuthGuard] }
    ]},
    { path: 'password', component: PasswordComponent, canActivate: [AuthGuard] },
    { path: 'users', component: UsersComponent, canActivate: [AuthGuard] },
    { path: 'signup', component: SignupComponent, canActivate: [AuthGuard] },
  ]},
  { path: 'auth', component: AuthComponent, children: [
    { path: 'signin', component: SigninComponent },
    { path: 'signup', component: SignupComponent }
  ]},
  { path: 'not-found', component: AppComponent, data: { message: "Page not found"} },
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
