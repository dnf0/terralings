import * as vscode from 'vscode';

/**
 * Manages the Terralings status bar item showing curriculum progress.
 */
export class TerralingsStatusBar implements vscode.Disposable {
  private readonly statusBarItem: vscode.StatusBarItem;

  constructor() {
    this.statusBarItem = vscode.window.createStatusBarItem(
      vscode.StatusBarAlignment.Left,
      100
    );
    this.statusBarItem.name = 'Terralings Progress';
    this.statusBarItem.command = 'workbench.view.extension.terralings-exercises';
    this.statusBarItem.tooltip = 'Click to view Curriculum & Exercises';
    this.update(0, 0);
  }

  /**
   * Updates the progress indicator with completed and total exercise counts.
   */
  public update(completed: number, total: number): void {
    const percentage = total > 0 ? Math.round((completed / total) * 100) : 0;
    this.statusBarItem.text = `$(zap) Terralings: ${completed}/${total} (${percentage}%)`;
    this.statusBarItem.tooltip = `Terralings Progress: ${completed} of ${total} completed (${percentage}%)\nClick to view Curriculum & Exercises`;
    this.statusBarItem.show();
  }

  /**
   * Explicitly hides the status bar item.
   */
  public hide(): void {
    this.statusBarItem.hide();
  }

  /**
   * Explicitly shows the status bar item.
   */
  public show(): void {
    this.statusBarItem.show();
  }

  /**
   * Disposes the underlying status bar item resource.
   */
  public dispose(): void {
    this.statusBarItem.dispose();
  }
}
