# GIT基本命令

1. 配置Git：

   在开始使用Git之前，您需要配置您的用户名和电子邮件地址。这些信息将用于识别您的提交。

   ```bash
   git config --global user.name "Your Name"
   git config --global user.email "your.email@example.com"
   ```

2. 克隆现有仓库：

   要从现有远程仓库创建本地副本，请使用以下命令：

   ```bash
   git clone https://github.com/username/repository.git
   ```

3. 初始化新仓库：

   要在现有项目中创建一个新的Git仓库，请使用以下命令：

   ```bash
   git init
   ```

4. 查看状态：

   要查看您的工作目录中已更改但尚未提交的文件，请使用以下命令：

   ```bash
   git status
   ```

5. 添加文件：

   要将新文件或已更改的文件添加到暂存区，请使用以下命令：

   ```bash
   git add <file>
   ```

   要一次添加多个文件或所有更改的文件，请使用：

   ```bash
   git add .
   ```

6. 提交更改：

   要将暂存区中的更改提交到仓库，请使用以下命令：

   ```bash
   git commit -m "Your commit message"
   ```

7. 查看日志：

   要查看仓库的提交日志，请使用以下命令：

   ```bash
   git log
   ```

8. 拉取更改：

   要从远程仓库获取最新更改并合并到当前分支，请使用以下命令：

   ```bash
   git pull origin main
   ```

9. 推送更改：

   要将本地分支的更改推送到远程仓库，请使用以下命令：

   ```bash
   git push origin main
   ```

10. 分支：

    要列出所有分支，请使用以下命令：

    ```bash
    git branch
    ```

    要创建一个新分支，请使用：

    ```bash
    git branch new-feature
    ```

    要删除一个分支，请使用：

    ```bash
    git branch -d new-feature
    ```

11. 切换分支：

    要切换到另一个分支，请使用以下命令：

    ```bash
    git checkout new-feature
    ```

12. 合并分支：

    要将一个分支的更改合并到当前分支，请使用以下命令：

    ```bash
    git merge new-feature
    ```