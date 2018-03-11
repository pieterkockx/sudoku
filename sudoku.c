#include <stdio.h>

int solve(size_t, int *);
void print(int *);

unsigned long long nsolutions = 0;

int solve(size_t start, int *grid)
{
    size_t i, j, n;
    int set[9] = {
      1, 1, 1, 1, 1, 1, 1, 1, 1,
    };

    for (n = 0; n < 81; n++)
        if (grid[n] == 0)
            break;

    if (n == 81) {
        nsolutions++;
        print(grid);
        return 1;
    }

    for (i = 9*(n / 9); i < 9*(n/9+1); i++)
        if (grid[i] > 0)
            set[grid[i] - 1] = 0;

    for (i = (n % 9); i < 81; i += 9)
        if (grid[i] > 0)
            set[grid[i] - 1] = 0;

    for (i = 27*(n / 27); i < 27*(n/27+1); i += 9)
        for (j = 3*((n % 9)/3); j < 3*((n % 9)/3 + 1); j++)
            if (grid[i + j] > 0)
                set[grid[i + j] - 1] = 0;

    for (i = start; i < 9; i++)
        if (set[i] == 1) {
            grid[n] = i + 1;
            if (solve(0, grid)) {
                grid[n] = 0;
                solve(i + 1, grid);
                return 1;
            }
            grid[n] = 0;
        }

    return 0;
}

int main(void)
{
    int grid[81] = {
          9, 0, 0, 0, 0, 0, 0, 0, 1,
          0, 0, 0, 0, 0, 3, 0, 8, 5,
          0, 0, 1, 0, 2, 0, 0, 0, 0,
          0, 0, 0, 5, 0, 7, 0, 0, 0,
          0, 0, 4, 0, 0, 0, 1, 0, 0,
          0, 9, 0, 0, 0, 0, 0, 0, 0,
          5, 0, 0, 0, 0, 0, 0, 7, 3,
          0, 0, 2, 0, 1, 0, 0, 0, 0,
          0, 0, 0, 0, 4, 0, 0, 0, 9,
    };

    if (solve(0, grid))
        printf("found %llu different solutions\n", nsolutions);
}

void print(int *grid)
{
    size_t i;
    for (i = 0; i < 81; i++) {
        printf("%d", grid[i]);
        if ( ((i + 1) % 9) == 0)
            printf("\n");
        else
            printf(" ");
    }
    printf("\n");
}
