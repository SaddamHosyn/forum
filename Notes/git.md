Branching Strategy
Create Feature Branches: Instead of working directly on the main branch, create separate branches for each feature or bug fix. This allows you to work in isolation without affecting the stable codebase. Use descriptive names for your branches, such as feature/login or bugfix/issue-12335.
Never Work on Main: The main branch should remain stable and clean. Always perform your development work on feature branches and merge them into the main branch only after thorough testing and code reviews24.
Keep Branches Short-Lived: Aim to keep your branches short-lived by merging them back into the main branch as soon as the feature is complete and tested. This helps avoid long-lived branches that can become outdated3.
Committing Changes
Commit Often: Make frequent commits to capture your progress. This not only protects your work but also allows you to track changes more effectively. Use meaningful commit messages that describe what changes were made15.
Use Pull Requests: Once you finish a feature, submit a pull request (PR) to merge your changes into the main branch. This allows other team members to review your code, suggest improvements, and ensure quality before integration45.
Review Process: Establish a practice where no one merges their own code without peer review. This encourages collaboration and helps catch potential issues early4.
Keeping Your Main Branch Clean
Regular Updates: Regularly pull updates from the main branch into your feature branches to minimize merge conflicts later. This ensures that your work is based on the most recent version of the codebase5.
Delete Merged Branches: After merging a feature branch back into the main branch, delete the branch both locally and on GitHub to keep the repository clean and organized25.
Use Git Flow or Similar Models: Consider adopting a branching model like Git Flow, which uses multiple primary branches (such as develop and main) to manage features and releases more effectively7.
By following these practices, you can maintain a clean main branch while fostering collaboration within your team, ultimately leading to a more successful project outcome.




Basic Git Branch Commands
List Local Branches:
git branch
This command displays all local branches in your repository.


List All Branches (Local and Remote):
git branch -a
This shows both local and remote branches.


Create a New Branch:
git checkout -b <branch-name>
This creates a new branch and switches to it immediately.

Switch to an Existing Branch:
git checkout <branch-name>
Use this command to switch to an existing branch.

Merge a Branch into Current Branch:
git merge <branch-name>
This merges the specified branch into the currently checked-out branch.

Delete a Local Branch:
git branch -d <branch-name>
This deletes the specified branch if it has been merged. Use -D to force delete it regardless of its merge status:
git branch -D <branch-name>
Push a New Branch to Remote:
git push -u origin <branch-name>
This pushes your new branch to the remote repository and sets up tracking.

Delete a Remote Branch:
git push origin --delete <branch-name>
This command removes the specified branch from the remote repository.

Check Current Branch:
git status
This shows the current branch you are on and any changes staged for commit.

Rebase Your Current Branch onto Another:
git rebase <target-branch>
