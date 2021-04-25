import { Project } from '@service/model';
import { BackendApiClient } from '@service/ApplicationServices';
import { StatService, StatServiceImpl } from '@service/stat';

export type PaymentPlan = {
  name: string,
  id: string,
  events_limit: number
}

export const paymentPlans: Record<string, PaymentPlan> = {
  free: { name: 'Startup (free)', id: 'free', events_limit: 250_000 },
  growth: { name: 'Growth', id: 'growth', events_limit: 1_000_000 },
  premium: { name: 'Premium', id: 'premium', events_limit: 10_000_000 },
  enterprise: { name: 'Enterprise', id: 'enterprise', events_limit: null }
}



/**
 * Status of current payment plan
 */
export class PaymentPlanStatus {
  private _currentPlan: PaymentPlan;

  private _eventsThisMonth: number;

  private _stat: StatService;

  public async init(project: Project, backendApiClient: BackendApiClient) {
    if (!project.planId) {
      this._currentPlan = paymentPlans.free;
    } else {
      this._currentPlan = paymentPlans[project.planId];
      if (!this._currentPlan) {
        throw new Error(`Unknown plan ${project.planId}`);
      }
    }
    this._stat = new StatServiceImpl(backendApiClient, project, true);
    var date = new Date();
    const stat = await this._stat.get(new Date(date.getFullYear(), date. getMonth(), 1), new Date(date.getFullYear(), date. getMonth() + 1, 0), 'day');
    this._eventsThisMonth = stat.reduce((res, item) => {
      res += item['data'];
      return res;
    }, 0);
  }

  get currentPlan(): PaymentPlan {
    return this._currentPlan;
  }

  get eventsThisMonth(): number {
    return this._eventsThisMonth;
  }

}