# Zammad Cleaner

Zammad Cleaner is a simple tool written in Go that interacts with the Zammad API to delete old tickets (Because I can not find any other tool which can do it).

## Usage

1. Download the latest release from the [Releases](https://github.com/martin-korf/zammad_cleaner/releases) page.

2. Navigate to the extracted directory.

3. Run the executable with the following command, replacing `YOUR_API_TOKEN` and `https://zammad.example.com` with your actual API token and Zammad base URL:

   ```bash
   ./zammad_cleaner --token=YOUR_API_TOKEN --base-url=https://zammad.example.com

Optionally, you can specify a cut-off date for ticket deletion using the --cut-off-date flag. The default cut-off date is 2021-07-01. The format for the date should be YYYY-MM-DD.

    ./zammad_cleaner --token=YOUR_API_TOKEN --base-url=https://zammad.example.com --cut-off-date=2021-07-01

## Description

Zammad Cleaner interacts with a Zammad helpdesk API to retrieve a list of tickets and delete those that were closed before a specified cut-off date. It uses the provided API token and base URL to authenticate and make API requests.

## Requirements

API token from your Zammad helpdesk provider
Access to the Zammad helpdesk API with permission to delete tickets.
