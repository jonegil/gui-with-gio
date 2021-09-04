
# Load analysis 
# Generate cpu logs from top, using
# First compile with wInvalidate
#   pgrep 11_improved_animation
#   top -l 0 -s 1  -pid 8849 -stats cpu | awk 'NR%13==0; fflush(stdout)' > wInvalidate.txt

# recompile top opInvalidateOp
#   pgrep 11_improved_animation
#   top -l 0 -s 1  -pid 9144 -stats cpu | awk 'NR%13==0; fflush(stdout)' > opInvalidateOpAdd.txt

# With compute renderer
#   GIORENDERER=forcecompute go run .
#   pgrep 11_improved_animation
#   top -l 0 -s 1  -pid 9144 -stats cpu | awk 'NR%13==0; fflush(stdout)' > opInvalidateOpAdd.txt

# Run Egg timer for 60 seconds
  library(data.table)
  library(ggplot2)

# Source
  d1 = fread("opInvalidateOpAdd.txt")
  d2 = fread("wInvalidate.txt")

# Tag
  d1[, method:="opInvalidateOpAdd"]
  d2[, method:="wInvalidate"]

# Combine and arrange
  d = rbind(d1,d2)
  setcolorder(d, c(1,3,2))
  d[, sec:=-2:(.N-3), method]
  md = melt(d, id.vars=c("method", "sec"))

# Visualize
  p1 = ggplot(md, aes(x=sec, y=value, color=method)) + 
      geom_line() +
      labs(x="Time (s)", 
           y="CPU load (%)", 
           title="Invalidation CPU load", 
           subtitle="Egg timer 60 second with opInvalidateOpAdd and wInvalidate", 
           caption="MacOS, MacBook Air (2017) 1.8 GHz Dual-Core Intel Core i5, Aug 8th 2021
           Same Macbook, GIORENDERER=forcecompute, Sep 4th 2021
           Windows 10, Intel i5 3570K Quad Core, Aug 14th 2021") +
      facet_wrap(~variable)
  
  p2 = ggplot(md[variable=="Windows" | (variable=="Macbook" & value > 10) | (variable=="Macbook_forcecompute" & value > 10)], aes(x=value, color=method)) + 
    geom_density() +
    labs(x="CPU load (%)", 
         title="Invalidation CPU load", 
         subtitle="Egg timer 60 second with opInvalidateOpAdd and wInvalidate", 
         caption="MacOS, MacBook Air (2017) 1.8 GHz Dual-Core Intel Core i5, Aug 8th 2021
         Same Macbook, GIORENDERER=forcecompute, Sep 4th 2021
         Windows 10, Intel i5 3570K Quad Core, Aug 14th 2021") +
    facet_wrap(~variable, scales="free_x")

# Save
  ggsave("../../11_invalidate_cpu_load.png", p1, width=8, height=4, dpi="print")
  ggsave("../../11_invalidate_cpu_density.png", p2, width=8, height=4, dpi="print")
